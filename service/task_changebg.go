package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"
	"wechat-dev/thinking/utils"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"github.com/redis/go-redis/v9"
)

type ChangeBGTask struct {
	TSender     *TextSender     `inject:""`
	ISender     *ImageSender    `inject:""`
	ElemService *ElementService `inject:""`
}

const (
	changeBgStateTodo  = "todo"
	changeBgStateDoing = "doing"
	changeBgStateFail  = "fail"
	changeBgExpire     = 1 * time.Minute
)

func (c *ChangeBGTask) Name() taskType {
	return taskTypeChangeBG
}

func (c *ChangeBGTask) Service(ctx *gin.Context, req models.RequestRawMessage) error {
	tt, ts, err := getCacheTask(ctx, req.FromUserName)
	if err != nil && err != redis.Nil {
		glg.Error("get task error", err, req)
		return err
	}

	if err == redis.Nil {
		return c.pre(ctx, req)
	}

	if tt != c.Name() {
		err = errors.New("task type is not changebg")
		glg.Error(err, req)
		return err
	}

	if req.Content == "退出" {
		return cleanCacheTask(ctx, req.FromUserName)
	}

	switch ts {
	case changeBgStateTodo:
		return c.start(ctx, req)
	case changeBgStateDoing:
		return c.check(ctx, req)
	case changeBgStateFail:
		return c.fail(ctx, req)
	default:
		mediaID := ts
		return c.exit(ctx, req, mediaID)
	}
}

// pre 点击按钮进入功能
func (c *ChangeBGTask) pre(ctx *gin.Context, req models.RequestRawMessage) error {
	err := setCacheTask(ctx, req.FromUserName, c.Name(), changeBgStateTodo, changeBgExpire)
	if err != nil {
		glg.Error("set cache task error", err, req)
		c.TSender.Send(ctx, req, "哎呀出错了！请点击按钮重试~")
		return err
	}

	glg.Info("enter change bg", req.FromUserName)

	c.TSender.Send(ctx, req, `请1分钟内上传一张图片
	· 要求尽量上传清晰的人物照片
	· 图片大小不要超过8M
	· 如果退出这个功能，请回复“退出”`)

	return nil
}

// start 开始上传图片，要求用户上传图片
func (c *ChangeBGTask) start(ctx *gin.Context, req models.RequestRawMessage) error {
	if req.MsgType != "image" {
		glg.Info("req not image", req)
		c.TSender.Send(ctx, req, "请上传图片或者回复“退出”~")
		return nil
	}

	if req.PicURL == "" {
		glg.Info("url is empty", req)
		c.TSender.Send(ctx, req, "请上传图片或者回复“退出”~")
		return nil
	}

	err := setCacheTask(ctx, req.FromUserName, c.Name(), changeBgStateDoing, changeBgExpire)
	if err != nil {
		glg.Error("set task error", err, req)
		return err
	}

	go func() {
		imgBytes, err := c.changeBG(ctx, req.PicURL)
		if err != nil {
			setCacheTask(ctx, req.FromUserName, c.Name(), changeBgStateFail, changeBgExpire)
			glg.Error("change bg error", err, req)
			return
		}
		filePath := fmt.Sprintf("./static/%s_%d.jpg", req.FromUserName, req.MsgID)
		err = utils.SaveImg(imgBytes, filePath)
		if err != nil {
			setCacheTask(ctx, req.FromUserName, c.Name(), changeBgStateFail, changeBgExpire)
			glg.Error("save img error", err, req)
			return
		}
		rsp, err := c.ElemService.Upload(ctx, filePath)
		if err != nil {
			setCacheTask(ctx, req.FromUserName, c.Name(), changeBgStateFail, changeBgExpire)
			glg.Error("upload img error", err, req)
			return
		}
		err = utils.DelImg(filePath)
		if err != nil {
			glg.Error("del img error", err, req)
		}
		err = setCacheTask(ctx, req.FromUserName, c.Name(), rsp.MediaID, changeBgExpire)
		if err != nil {
			setCacheTask(ctx, req.FromUserName, c.Name(), changeBgStateFail, changeBgExpire)
			glg.Error("set task done error", err, req, rsp)
			return
		}
	}()

	glg.Info("start change bg", req.FromUserName)

	c.TSender.Send(ctx, req, `图片正在处理中，请稍后......
			Tip: 我不能主动给你发图片处理结果。所以，请回复“1”获取处理后的图片！`)

	return nil
}

// check 图片正在处理中，用户检查图片处理状态，要求用户输入“1”
func (c *ChangeBGTask) check(ctx *gin.Context, req models.RequestRawMessage) error {
	if req.MsgType != "text" || req.Content != "1" {
		utils.NoNeedResponse(ctx)
		return nil
	}

	glg.Info("enter check", req.FromUserName)

	c.TSender.Send(ctx, req, `图片处理中，请稍后......
					稍等几秒后，再回复“1”获取`)
	return nil
}

// fail 图片处理失败
func (c *ChangeBGTask) fail(ctx *gin.Context, req models.RequestRawMessage) error {
	if req.MsgType != "text" || req.Content != "1" {
		utils.NoNeedResponse(ctx)
		return nil
	}

	glg.Info("enter fail", req.FromUserName)

	cleanCacheTask(ctx, req.FromUserName)

	c.TSender.Send(ctx, req, "不好意思！图片处理失败了~")
	return nil
}

// exit 图片处理完，用户获取图片
func (c *ChangeBGTask) exit(ctx *gin.Context, req models.RequestRawMessage, mediaID string) error {
	if req.MsgType != "text" || req.Content != "1" {
		utils.NoNeedResponse(ctx)
		return nil
	}
	cleanCacheTask(ctx, req.FromUserName)
	glg.Info("enter exit", req.FromUserName)
	c.ISender.Send(ctx, req, mediaID)
	return nil
}

// changeBG 调用外部接口处理图片
func (c *ChangeBGTask) changeBG(ctx context.Context, picUrl string) ([]byte, error) {
	var (
		url = "https://api.remove.bg/v1.0/removebg"
	)
	// 创建 multipart 上传对象
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加其他字段
	writer.WriteField("size", "preview")
	writer.WriteField("image_file_b64", "")
	writer.WriteField("semitransparency", "true")
	writer.WriteField("position", "original")
	writer.WriteField("bg_color", "red") // TODO 背景颜色
	writer.WriteField("scale", "original")
	writer.WriteField("image_url", picUrl)
	writer.WriteField("roi", "0% 0% 100% 100%")
	writer.WriteField("crop", "false")
	writer.WriteField("channels", "rgba")
	writer.WriteField("bg_image_url", "")
	writer.WriteField("format", "auto")
	writer.WriteField("bg_image_file", "")
	writer.WriteField("type", "auto")
	writer.WriteField("crop_margin", "0")
	writer.WriteField("add_shadow", "false")
	writer.WriteField("type_level", "1")

	// 写入 multipart 结束标识
	writer.Close()

	// 发送上传请求
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		glg.Errorf("create request error: %v. picUrl=%s\n", err, picUrl)
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "image/*")
	req.Header.Set("X-API-Key", config.ChangeBgApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("send request error: %v. picUrl=%s\n", err, picUrl)
		return nil, err
	}

	// 解析结果
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("read response error: %v. picUrl=%s\n", err, picUrl)
		return nil, err
	}

	var errorRsp models.ResponseChangeBG
	err = json.Unmarshal(respBytes, &errorRsp)
	if err == nil {
		err = errors.New("err code != 0")
		glg.Errorf("err=%v\n", errorRsp)
		return nil, err
	}

	glg.Infof("change bg exec doner. esponse header: used credit=%s\n", resp.Header.Get("X-Credits-Charged"))

	return respBytes, nil
}

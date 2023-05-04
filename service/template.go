// Description: 使用微信模板消息功能，定时提醒用户
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// SetTemplate 设置待处理模板消息
func SetTemplate(ctx *gin.Context, second int, msg models.RequestRawMessage) {
	// notice after second * 10
	noticeTime := float64(time.Now().Unix() + int64(second*10))
	err := config.GetRedisClient().ZAdd(ctx, config.RedisMemberKey, redis.Z{
		Score:  noticeTime,
		Member: msg.FromUserName,
	}).Err()
	if err != nil {
		glg.Errorf("Zadd error. err=%v\n", err)
		ctx.String(http.StatusOK, "success")
		return
	}

	// 需要根节点是xml
	type xml struct {
		models.RequestRawMessage
	}
	ctx.XML(http.StatusOK, xml{
		RequestRawMessage: models.RequestRawMessage{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			MsgType:      msg.MsgType,
			CreateTime:   time.Now().Unix(),
			Content:      fmt.Sprintf("你将于%d秒后收到提醒！", second*10),
		},
	})
}

// DoTemplate 处理模板消息
func DoTemplate() {
	go func() {
		for {
			notice()
			time.Sleep(500 * time.Millisecond)
		}
	}()
	glg.Info("start to send template ......")
}

func notice() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	now := time.Now().Unix()
	result, err := config.GetRedisClient().ZRangeByScore(ctx, config.RedisMemberKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(now, 10),
	}).Result()
	if err != nil {
		glg.Errorf("get notice memeber list error. err=%v\n", err)
		return
	}
	if len(result) == 0 {
		return
	}

	var ts = &TokenService{}
	accessToken, err := ts.GetAccessToken(ctx)
	if err != nil {
		glg.Errorf("get access token error. err=%v\n", err)
		return
	}

	for _, each := range result {
		// each is openID
		_ = send(ctx, each, accessToken.Token)
	}
	err = config.GetRedisClient().ZRem(ctx, config.RedisMemberKey, result).Err()
	if err != nil {
		glg.Infof("del memebers error. err=%v\n", err)
		return
	}
	glg.Infof("send %d templates\n", len(result))
}

func send(ctx context.Context, toUser, token string) error {
	var (
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", token)
	)

	templateMsg := models.NewTemplateMsg(toUser)
	body, err := json.Marshal(templateMsg)
	if err != nil {
		glg.Errorf("marshal error. err=%v\n", err)
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		glg.Errorf("new request error. err=%v\n", err)
		return err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		glg.Errorf("do http request error. err=%v\n", err)
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		glg.Errorf("status is not ok error. err=%v\n", err)
		return errors.New("status is not ok")
	}

	body, err = io.ReadAll(rsp.Body)
	if err != nil {
		glg.Errorf("read body error. err=%v\n", err)
		return err
	}

	var commonRsp models.CommonResponse
	err = json.Unmarshal(body, &commonRsp)
	if err != nil {
		glg.Errorf("unmarshal error. err=%v\n", err)
		return err
	}

	return nil
}

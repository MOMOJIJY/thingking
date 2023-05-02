package handler

import (
	"strconv"
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"

	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct{}

// Verify 开发者身份验证
func (h *Handler) Verify(ctx *gin.Context) {
	// 1. 获取参数
	signature, _ := ctx.GetQuery("signature")
	timestamp, _ := ctx.GetQuery("timestamp")
	nonce, _ := ctx.GetQuery("nonce")
	echostr, _ := ctx.GetQuery("echostr")
	if signature == "" || timestamp == "" || nonce == "" || echostr == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid parameters"))
	}

	// 2. 验证
	params := []string{config.Token, timestamp, nonce}
	sort.Strings(params)
	hashcode := hash(params)

	log.Printf("Request: signature=%s, timestamp=%s, nonce=%s, echostr=%s, hashcode=%s\n",
		signature, timestamp, nonce, echostr, hashcode)

	if hashcode == signature {
		ctx.String(http.StatusOK, echostr)
	}
	ctx.String(http.StatusOK, "")
}

// HandleMsgFromUser 获取用户消息
func (h *Handler) HandleMsgFromUser(ctx *gin.Context) {
	var msg models.MsgFromUser
	err := ctx.ShouldBindXML(&msg)
	if err != nil {
		log.Printf("bind json error. err=%v\n", err)
		// 如果不想重试需要回复"success"
		ctx.String(http.StatusOK, "success")
		return
	}
	log.Printf("from=%s, to=%s, ctime=%d, msgId=%d, msgType=%s, content=%s\n", msg.FromUserName,
		msg.ToUserName, msg.CreateTime, msg.MsgId, msg.MsgType, msg.Content)

	if msg.MsgType == "text" {
		h.replyTextMsg(ctx, msg)
	} else if msg.MsgType == "event" {
		h.replyEventMsg(ctx, msg)
	} else {
		// 如果不想重试需要回复"success"
		ctx.String(http.StatusOK, "success")
		return
	}

}

func (h *Handler) replyTextMsg(ctx *gin.Context, msg models.MsgFromUser) {
	if second, err := strconv.Atoi(msg.Content); err == nil {
		h.doTemplate(ctx, second, msg)
		return
	}

	// 需要根节点是xml
	type xml struct {
		models.MsgFromUser
	}
	ctx.XML(http.StatusOK, xml{
		MsgFromUser: models.MsgFromUser{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			MsgType:      msg.MsgType,
			CreateTime:   time.Now().Unix(),
			Content:      "已收到！",
		},
	})
}

func (h *Handler) doTemplate(ctx *gin.Context, second int, msg models.MsgFromUser) {
	// notice after second * 10
	noticeTime := float64(time.Now().Unix() + int64(second*10))
	err := config.GetRedisClient().ZAdd(ctx, config.RedisMemberKey, redis.Z{
		Score:  noticeTime,
		Member: msg.FromUserName,
	}).Err()
	if err != nil {
		log.Printf("Zadd error. err=%v\n", err)
		ctx.String(http.StatusOK, "success")
		return
	}

	// 需要根节点是xml
	type xml struct {
		models.MsgFromUser
	}
	ctx.XML(http.StatusOK, xml{
		MsgFromUser: models.MsgFromUser{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			MsgType:      msg.MsgType,
			CreateTime:   time.Now().Unix(),
			Content:      fmt.Sprintf("你将于%d秒后收到提醒！", second*10),
		},
	})
}

func (h *Handler) replyEventMsg(ctx *gin.Context, msg models.MsgFromUser) {
	if msg.Event != "CLICK" {
		ctx.String(http.StatusOK, "success")
		return
	}

	// 需要根节点是xml
	type xml struct {
		models.MsgFromUser
	}

	if msg.EventKey == "V1001_FOLLOW" {
		ctx.XML(http.StatusOK, xml{
			MsgFromUser: models.MsgFromUser{
				ToUserName:   msg.FromUserName,
				FromUserName: msg.ToUserName,
				MsgType:      "text",
				CreateTime:   time.Now().Unix(),
				Content:      "转发成功！",
			},
		})
	} else {
		ctx.XML(http.StatusOK, xml{
			MsgFromUser: models.MsgFromUser{
				ToUserName:   msg.FromUserName,
				FromUserName: msg.ToUserName,
				MsgType:      "text",
				CreateTime:   time.Now().Unix(),
				Content:      "点赞、投币、收藏，一键三连！",
			},
		})
	}
}

// hash get hash code
func hash(params []string) string {
	var s string
	for _, each := range params {
		s += each
	}
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

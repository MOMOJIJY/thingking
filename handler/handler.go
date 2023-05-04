package handler

import (
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"
	"wechat-dev/thinking/service"
	"wechat-dev/thinking/utils"

	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

type Handler struct {
}

// Verify 开发者身份验证
func (h *Handler) Verify(ctx *gin.Context) {
	// 1. 获取参数
	signature, _ := ctx.GetQuery("signature")
	timestamp, _ := ctx.GetQuery("timestamp")
	nonce, _ := ctx.GetQuery("nonce")
	echostr, _ := ctx.GetQuery("echostr")
	if signature == "" || timestamp == "" || nonce == "" || echostr == "" {
		err := errors.New("invalid parameters")
		glg.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	// 2. 验证
	params := []string{config.Token, timestamp, nonce}
	sort.Strings(params)
	hashcode := hash(params)

	glg.Infof("Request: signature=%s, timestamp=%s, nonce=%s, echostr=%s, hashcode=%s\n",
		signature, timestamp, nonce, echostr, hashcode)

	if hashcode == signature {
		ctx.String(http.StatusOK, echostr)
	}
	ctx.String(http.StatusOK, "")
}

// HandleMsg 处理用户信息
func (h *Handler) HandleMsg(ctx *gin.Context) {
	var msg models.RequestRawMessage
	err := ctx.ShouldBindXML(&msg)
	if err != nil {
		glg.Errorf("bind xml error. err=%v\n", err)
		utils.NoNeedResponse(ctx)
		return
	}
	glg.Info("receive user msg", msg)

	task, err := service.Prepare(ctx, msg)
	if err != nil {
		glg.Error("prepare error", err, msg)
		utils.NoNeedResponse(ctx)
		return
	}

	err = task.Service(ctx, msg)
	if err != nil {
		glg.Error("service error", err, msg)
		utils.NoNeedResponse(ctx)
		return
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

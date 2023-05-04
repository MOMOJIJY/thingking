package service

import (
	"encoding/xml"
	"net/http"
	"time"
	"wechat-dev/thinking/models"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

type TextSender struct{}

func (s *TextSender) Send(ctx *gin.Context, req models.RequestRawMessage, content string) error {
	textMsg := models.ResponseTextMessage{}
	textMsg.MsgType = models.MsgTypeText
	textMsg.FromUserName = req.ToUserName
	textMsg.ToUserName = req.FromUserName
	textMsg.CreateTime = time.Now().Unix()
	textMsg.Content = content

	responseBody, err := xml.Marshal(textMsg)
	if err != nil {
		glg.Error("xml marshal error", err, textMsg)
		return err
	}

	ctx.Data(http.StatusOK, models.XmlContentType, responseBody)

	glg.Info("Send text msg", textMsg)
	return nil
}

package service

import (
	"encoding/xml"
	"net/http"
	"time"
	"wechat-dev/thinking/models"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

type ImageSender struct{}

func (s *ImageSender) Send(ctx *gin.Context, req models.RequestRawMessage, mediaID string) error {
	imageMsg := models.ResponseImageMessage{}
	imageMsg.MsgType = models.MsgTypeImage
	imageMsg.FromUserName = req.ToUserName
	imageMsg.ToUserName = req.FromUserName
	imageMsg.CreateTime = time.Now().Unix()
	imageMsg.MediaID = mediaID

	responseBody, err := xml.Marshal(imageMsg)
	if err != nil {
		glg.Error("xml marshal error", err, imageMsg)
		return err
	}

	ctx.Data(http.StatusOK, models.XmlContentType, responseBody)

	glg.Info("Send text msg", imageMsg)
	return nil
}

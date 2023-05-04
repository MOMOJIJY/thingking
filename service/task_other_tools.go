package service

import (
	"wechat-dev/thinking/models"

	"github.com/gin-gonic/gin"
)

type OtherToolTask struct {
	TSender *TextSender `inject:""`
}

func (o *OtherToolTask) Name() taskType {
	return "other"
}

func (o *OtherToolTask) Service(ctx *gin.Context, req models.RequestRawMessage) error {
	o.TSender.Send(ctx, req, `更多精彩，陆续上线中：
				(1) AI智能律师：帮忙提供日常法律援助；
				(2) 图片增强：优化模糊图片。`)
	return nil
}

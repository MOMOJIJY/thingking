package service

import (
	"wechat-dev/thinking/models"

	"github.com/gin-gonic/gin"
)

type DefaultTask struct {
	TSender *TextSender `inject:""`
}

func (t *DefaultTask) Name() taskType {
	return "default"
}

func (t *DefaultTask) Service(ctx *gin.Context, req models.RequestRawMessage) error {
	t.TSender.Send(ctx, req, "点下面的菜单选择功能哦~")
	return nil
}

package service

import (
	"wechat-dev/thinking/models"

	"github.com/gin-gonic/gin"
)

type ThreeTask struct {
	TSender *TextSender `inject:""`
}

func (t *ThreeTask) Name() taskType {
	return "three"
}

func (t *ThreeTask) Service(ctx *gin.Context, req models.RequestRawMessage) error {
	t.TSender.Send(ctx, req, "投币点赞收藏，一键三连！")
	return nil
}

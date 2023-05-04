package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"

	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"github.com/redis/go-redis/v9"
)

type Task interface {
	Name() taskType
	Service(*gin.Context, models.RequestRawMessage) error
}

type taskType string

const (
	taskKey   = "task_%s" // task_{openID}
	taskValue = "%s_%s"   // {taskType}_{taskState}

	taskTypeDefault    = "default"
	taskTypeChangeBG   = "changebg"
	taskTypeOtherTools = "other"
	taskTypeThreeLian  = "three"
)

var taskCenter map[taskType]Task

func init() {
	taskCenter = make(map[taskType]Task)
	g := inject.Graph{}
	var (
		changebgTask  = new(ChangeBGTask)
		otherToolTask = new(OtherToolTask)
		threeTask     = new(ThreeTask)
	)
	err := g.Provide(
		&inject.Object{Value: changebgTask},
		&inject.Object{Value: otherToolTask},
		&inject.Object{Value: threeTask},
	)
	if err != nil {
		glg.Fatal(err)
	}
	err = g.Populate()
	if err != nil {
		glg.Fatal(err)
	}

	taskCenter[changebgTask.Name()] = changebgTask
	taskCenter[otherToolTask.Name()] = otherToolTask
	taskCenter[threeTask.Name()] = threeTask
}

// Prepare 处理消息前，根据消息和状态选择Task去处理
// @req: 用户发送的消息内容
func Prepare(ctx context.Context, req models.RequestRawMessage) (Task, error) {
	tt, _, err := getCacheTask(ctx, req.FromUserName)
	if err != nil && err != redis.Nil {
		glg.Error("get task error", err, req)
		return nil, err
	}

	// not task
	if err == redis.Nil {
		switch req.MsgType {
		case "event":
			return switchEventTask(req)
		default:
			task := new(DefaultTask)
			return task, nil
		}
	}

	// user has one task
	task, ok := taskCenter[tt]
	if !ok {
		glg.Warn("task not found", tt, req)
		task = new(DefaultTask)
		return task, nil
	}
	return task, nil
}

// switchEvent 根据事件选择任务
func switchEventTask(req models.RequestRawMessage) (Task, error) {
	switch req.EventKey {
	case models.EventKeyThreeLian:
		task := new(ThreeTask)
		return task, nil
	case models.EventKeyChangeBG:
		task := new(ChangeBGTask)
		return task, nil
	case models.EventKeyOtherTools:
		task := new(OtherToolTask)
		return task, nil
	default:
		err := errors.New("unknown event key")
		return nil, err
	}
}

// getCacheTask 从redis中获取任务
// @openID: 用户标识符
// @tt: task type
// @ts: task state
func getCacheTask(ctx context.Context, openID string) (tt taskType, ts string, err error) {
	val, err := config.GetRedisClient().Get(ctx, getTaskKey(openID)).Result()
	if err != nil {
		return
	}

	vals := strings.SplitN(val, "_", 2)
	tt, ts = taskType(vals[0]), vals[1]
	return
}

// setCacheTask 设置任务到redis
// @openID: 用户标识符
// @tt: task type
// @ts: task state
// @te: task expire
func setCacheTask(ctx context.Context, openID string, tt taskType, ts string, te time.Duration) error {
	return config.GetRedisClient().Set(ctx, getTaskKey(openID), getTaskValue(tt, ts), te).Err()
}

// cleanCacheTask 从redis中删除任务
// @openID: 用户标识符
func cleanCacheTask(ctx context.Context, openID string) error {
	return config.GetRedisClient().Del(ctx, getTaskKey(openID)).Err()
}

// getTaskKey 获取用户的任务的键
// @openID: 用户标识符
func getTaskKey(openID string) string {
	return fmt.Sprintf(taskKey, openID)
}

// getTaskKey 获取用户的任务的值
// @taskType: 任务类型
// @state: 任务状态
func getTaskValue(typ taskType, state string) string {
	return fmt.Sprintf(taskValue, typ, state)
}

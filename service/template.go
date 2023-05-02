// Description: 使用微信模板消息功能，定时提醒用户
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func DoTemplate() {
	go func() {
		for {
			notice()
			time.Sleep(500 * time.Millisecond)
		}
	}()
	log.Println("start to send template ......")
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
		log.Printf("get notice memeber list error. err=%v\n", err)
		return
	}
	if len(result) == 0 {
		return
	}

	var ts = &TokenService{}
	accessToken, err := ts.GetAccessToken(ctx)
	if err != nil {
		log.Printf("get access token error. err=%v\n", err)
		return
	}

	for _, each := range result {
		// each is openID
		_ = send(ctx, each, accessToken.Token)
	}
	err = config.GetRedisClient().ZRem(ctx, config.RedisMemberKey, result).Err()
	if err != nil {
		log.Printf("del memebers error. err=%v\n", err)
		return
	}
	log.Printf("send %d templates\n", len(result))
}

func send(ctx context.Context, toUser, token string) error {
	var (
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", token)
	)

	templateMsg := models.NewTemplateMsg(toUser)
	body, err := json.Marshal(templateMsg)
	if err != nil {
		log.Printf("marshal error. err=%v\n", err)
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("new request error. err=%v\n", err)
		return err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("do http request error. err=%v\n", err)
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		log.Printf("status is not ok error. err=%v\n", err)
		return errors.New("status is not ok")
	}

	body, err = io.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("read body error. err=%v\n", err)
		return err
	}

	var commonRsp models.CommonResponse
	err = json.Unmarshal(body, &commonRsp)
	if err != nil {
		log.Printf("unmarshal error. err=%v\n", err)
		return err
	}

	return nil
}

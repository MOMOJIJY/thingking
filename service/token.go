// Description: 从微信服务器获取access token，存到redis里面，后面直接从redis中获取
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"

	"github.com/kpango/glg"
	"github.com/redis/go-redis/v9"
)

type TokenService struct{}

func (t *TokenService) GetAccessToken(ctx context.Context) (models.AccessToken, error) {
	accessToken, err := t.loadAccessToken(ctx)
	if err != nil {
		return accessToken, err
	}

	if accessToken.Token != "" {
		return accessToken, nil
	}

	accessToken, err = t.realGetAccessToken()
	if err != nil {
		return accessToken, err
	}

	if accessToken.Token == "" {
		return accessToken, errors.New("get access token failed")
	}

	err = t.setAccessToken(ctx, accessToken)
	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

// realGetAccessToken 从微信服务器获取access_token
func (t *TokenService) realGetAccessToken() (models.AccessToken, error) {
	var accessToken models.AccessToken

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		config.AppID, config.AppSecret)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		glg.Errorf("New Request error. err=%v\n", err)
		return accessToken, err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		glg.Errorf("Do http request error. err=%v\n", err)
		return accessToken, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		glg.Errorf("code is not equal 200. err=%v\n", err)
		return accessToken, err
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		glg.Errorf("read body error. err=%v\n", err)
		return accessToken, err
	}

	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		glg.Errorf("json unmarshal error. err=%v\n", err)
		return accessToken, err
	}

	return accessToken, nil
}

// loadAccessToken 从redis读access_token
func (t *TokenService) loadAccessToken(ctx context.Context) (models.AccessToken, error) {
	var accessToken models.AccessToken
	val, err := config.GetRedisClient().Get(ctx, config.RedisTokenKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return accessToken, nil
		}
		return accessToken, err
	}

	err = json.Unmarshal(val, &accessToken)
	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

// setAccessToken 写access_token
func (t *TokenService) setAccessToken(ctx context.Context, accessToken models.AccessToken) error {
	val, err := json.Marshal(accessToken)
	if err != nil {
		return err
	}

	err = config.GetRedisClient().Set(ctx, config.RedisTokenKey, string(val), config.RedisTokenExpire).Err()
	if err != nil {
		return err
	}

	return nil
}

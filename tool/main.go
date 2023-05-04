package main

import (
	"context"
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"
	"wechat-dev/thinking/service"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kpango/glg"
)

func main() {
	config.InitRedis()
	config.InitBase()

	tk := &service.TokenService{}
	ctx := context.Background()

	accessToken, err := tk.GetAccessToken(ctx)
	if err != nil {
		glg.Fatal(err)
	}
	createMenu(accessToken.Token)
}

func createMenu(accessToken string) {
	menu := models.NewMenu()
	body, err := json.Marshal(menu)
	if err != nil {
		glg.Fatal(err)
		return
	}
	log.Println(string(body))

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", accessToken)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		glg.Fatal(err)
		return
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		glg.Fatal(err)
		return
	}
	defer rsp.Body.Close()

	body, err = io.ReadAll(rsp.Body)
	if err != nil {
		glg.Fatal(err)
		return
	}

	var commonResp models.CommonResponse
	err = json.Unmarshal(body, &commonResp)
	if err != nil {
		glg.Fatal(err)
		return
	}

	glg.Infof("errcode=%d, errmsg=%s\n", commonResp.ErrCode, commonResp.ErrMsg)
}

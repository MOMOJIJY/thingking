package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"wechat-dev/thinking/models"

	"github.com/kpango/glg"
)

const (
	uploadUrl = "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=image"
)

type ElementService struct{}

func (e *ElementService) Upload(ctx context.Context, filePath string) (models.ElementResponse, error) {
	var elemRsp models.ElementResponse
	var tk = &TokenService{}
	accessToken, err := tk.GetAccessToken(ctx)
	if err != nil {
		glg.Error("get access token error", err, filePath)
		return elemRsp, err
	}

	fileBytes, err := ioutil.ReadFile(filePath) // 读写方式打开
	if err != nil {
		glg.Error("read file error", err, filePath)
		return elemRsp, err
	}

	fileName := filepath.Base(filePath)
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("media", fileName)
	fileWriter.Write(fileBytes)
	contentType := bodyWriter.FormDataContentType() //contentType
	bodyWriter.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf(uploadUrl, accessToken.Token), bodyBuffer)
	if err != nil {
		glg.Error("http new request error", err, filePath, bodyBuffer.String())
		return elemRsp, err
	}
	req.Header.Set("Content-Type", contentType)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		glg.Error("do http error", err, filePath)
		return elemRsp, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		err = errors.New("upload response status is not ok")
		glg.Error(err, filePath)
		return elemRsp, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		glg.Error("read response body error", err, filePath, rsp)
		return elemRsp, err
	}

	var commonRsp models.CommonResponse
	err = json.Unmarshal(body, &commonRsp)
	if err != nil {
		glg.Error("json unmarshal for response body error", err, filePath, string(body))
		return elemRsp, err
	}
	if commonRsp.ErrCode != 0 {
		err = errors.New("upload response error")
		glg.Error(err, filePath)
		return elemRsp, err
	}

	err = json.Unmarshal(body, &elemRsp)
	if err != nil {
		glg.Error("json unmarshal for response body error", err, filePath, string(body))
		return elemRsp, err
	}

	glg.Info("upload element", filePath)
	return elemRsp, nil
}

func (e *ElementService) Get() {}

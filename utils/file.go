package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"

	"github.com/kpango/glg"
)

func SaveImg(body []byte, filePath string) error {
	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		glg.Error("decode image error", err, filePath)
		return err
	}

	// 保存为新图片文件
	outFile, err := os.Create(filePath)
	if err != nil {
		glg.Error("create file error", err, filePath)
		return err
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80})
	if err != nil {
		glg.Error("encode image error", err, filePath)
		return err
	}

	glg.Info("save success", filePath)
	return nil
}

func DelImg(filePath string) error {
	return nil
}

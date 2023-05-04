package config

import (
	"github.com/kpango/glg"
)

func InitLogger() {
	infoFile := glg.FileWriter("./static/info.log", 0666)
	errorFile := glg.FileWriter("./static/error.log", 0666)

	glg.Get().SetMode(glg.BOTH).
		SetLineTraceMode(glg.TraceLineShort).
		SetLevelWriter(glg.INFO, infoFile).
		SetLevelWriter(glg.DEBG, infoFile).
		SetLevelWriter(glg.ERR, errorFile).
		SetLevel(glg.DEBG).
		SetLevelColor(glg.INFO, glg.Green).
		SetLevelColor(glg.DEBG, glg.White).
		SetLevelColor(glg.ERR, glg.Red).
		EnableJSON()

	glg.Info("init logger")

}

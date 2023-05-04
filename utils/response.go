package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NoNeedResponse 不需要回复，回复“success”
func NoNeedResponse(ctx *gin.Context) {
	ctx.String(http.StatusOK, "success")
}

package service

import (
	"context"
	"testing"
	"wechat-dev/thinking/config"

	"github.com/stretchr/testify/require"
)

func TestGetAccessToken(t *testing.T) {
	config.InitRedis()

	var (
		ts  = &TokenService{}
		ctx = context.Background()
	)

	accessToken, err := ts.GetAccessToken(ctx)
	require.Nil(t, err)
	require.NotEmpty(t, accessToken)
}

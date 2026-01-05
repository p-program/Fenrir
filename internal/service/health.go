package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p-program/Fenrir/internal/core/config"
	"github.com/p-program/Fenrir/internal/core/logprovider"
	"github.com/p-program/Fenrir/internal/core/webprovider"
)

func NewHealthService(gin webprovider.MyGinEngine, l logprovider.Logger,
	config config.Config) HealthService {
	return HealthService{
		gin:    gin,
		l:      l,
		config: config,
	}
}

type HealthService struct {
	gin    webprovider.MyGinEngine
	l      logprovider.Logger
	config config.Config
}

func (s HealthService) Check(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusOK, struct {
		Code int `json:"code"`
	}{200})
}

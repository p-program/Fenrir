package core

import (
	"go.uber.org/fx"
	"github.com/p-program/Fenrir/internal/core/config"
	"github.com/p-program/Fenrir/internal/core/logprovider"
	"github.com/p-program/Fenrir/internal/core/webprovider"
)

var CoreModule = fx.Options(
	fx.Provide(config.NewFileConfig),
	fx.Provide(logprovider.GetLogger),
	//todo 集成数据库
	// fx.Provide(NewDatabase),
	fx.Provide(webprovider.NewGinEngine),
)

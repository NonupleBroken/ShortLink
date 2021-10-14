package server

import (
	"ShortLink/config"
	"ShortLink/db/redis"
	"ShortLink/logger"
	"ShortLink/web"
	"sync"
)

func WaitStop(wg *sync.WaitGroup) {
	defer func() {
		_ = logger.S.Sync()
		_ = logger.L.Sync()
	}()
	wg.Wait()
	logger.L.Info("server stop")
}

func Stop() {
	_ = redis.Client.Close()
}

func Run(cfg *config.Config) {
	err := logger.InitLogger(&cfg.Log)
	if err != nil {
		return
	}
	logger.L.Info("init logger success")

	err = redis.InitRedis(&cfg.Redis)
	if err != nil {
		logger.S.Error("init redis failed: ", err)
		return
	}
	logger.L.Info("init redis success")

	err = web.Start(cfg)
	if err != nil {
		logger.S.Error("web server start failed: ", err)
		return
	}
	logger.L.Info("web server start")
}

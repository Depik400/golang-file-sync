package app

import (
	"fmt"
	"golang-file-sync/internal/config"
	"golang-file-sync/internal/container"
)

func Run(configPath string) int {
	conf, _ := config.NewConfig(configPath)
	container.InitContainer(conf)
	container.GetContainer().Logger.Info(fmt.Sprintf("Server started with config:%v", conf)) //TODO READABLE CONFIG
	container.GetContainer().Delivery.Run()
	container.GetContainer().SyncService.Run()
	container.GetContainer().WatcherService.Run()
	<-make(chan interface{}) //TODO Сделать нормальный выход
	return 0
}

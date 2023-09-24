package app

import (
	"1._file-sync/internal/config"
	"1._file-sync/internal/container"
	"fmt"
)

func Run(configPath string) int {
	conf, _ := config.NewConfig(configPath)
	container.InitContainer(conf)
	container.GetContainer().Logger.Info(fmt.Sprintf("Server started with config:%v", conf)) //TODO READABLE CONFIG
	container.GetContainer().WatcherService.Run()
	<-make(chan interface{}) //TODO Сделать нормальный выход
	return 0
}

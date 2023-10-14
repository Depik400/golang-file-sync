package services

import (
	"fmt"
	"golang-file-sync/internal/delivery"
	"golang-file-sync/pkg/logger"
)

type SyncService struct {
	delivery delivery.IDelivery
	logger   logger.ILogger
}

func (service *SyncService) Run() {
	channel := service.delivery.GetMessageChannel()
	go func() {
		select {
		case message := <-channel:
			go service.logger.Info(fmt.Sprintf("message %s", message))
			go service.delivery.WriteMessage([]byte("Hello from server\n"), message.Connection)
			switch message.Cmd {
			case "file":
				{
				}
			}
		}
	}()
}

func NewSyncService(delivery delivery.IDelivery, iLogger logger.ILogger) *SyncService {
	return &SyncService{
		delivery,
		iLogger,
	}
}

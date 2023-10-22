package services

import (
	"fmt"
	"golang-file-sync/internal/delivery"
	"golang-file-sync/pkg/logger"
	"io"
	"log"
	"os"
	"path/filepath"
)

type SyncService struct {
	delivery delivery.IDelivery
	logger   logger.ILogger
	watcher  IWatchService
}

func (service *SyncService) Run() {
	channel := service.delivery.GetMessageChannel()
	go func() {
		for {
			message := <-channel
			go service.logger.Info(fmt.Sprintf("message %s", message))
			go service.delivery.WriteMessage([]byte("Hello from server\n"), message.Connection)
			switch message.Cmd {
			case "file":
				{
					buf := make([]byte, 1024*1024)
					for {
						n, err := message.Connection.Read(buf)
						if err != nil {
							if err != io.EOF {
								log.Println(err)
							}
							break
						}
						log.Printf("received: %s", string(buf[:n]))
						log.Printf("bytes: %d", n)
					}
					message.Connection.Close()
				}
			}
		}
	}()

	go func() {
		for {
			message := <-service.watcher.GetChannel()
			sockets := service.delivery.PushFileToAll(&delivery.IMessage{
				Cmd: "file",
				Metadata: delivery.IMetaData{
					FileName:     message.Filename,
					LocationName: message.DirectoryAlias,
					Action:       message.Op,
				},
			})
			for _, socket := range *sockets {
				file, err := os.Open(filepath.Join(message.Directory, message.Filename))
				if err != nil {
					log.Fatal(err)
				}
				pr, pw := io.Pipe()
				go func() {
					n, err := io.Copy(pw, file)
					if err != nil {
						log.Fatal(err)
					}
					pw.Close()
					log.Printf("copied to piped writer via the compressed writer: %d", n)
				}()

				n, err := io.Copy(socket, pr)
				if err != nil {
					log.Fatal(err)
				}
				socket.Close()
				log.Printf("copied to connection: %d", n)
			}
		}
	}()
}

func NewSyncService(delivery delivery.IDelivery, iLogger logger.ILogger, service IWatchService) *SyncService {
	return &SyncService{
		delivery,
		iLogger,
		service,
	}
}

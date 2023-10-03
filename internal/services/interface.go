package services

type IWatchService interface {
	Run()
	Stop()
}

type ISyncService interface {
	Run()
}

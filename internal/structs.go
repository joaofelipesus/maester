package internal

type Config struct {
	ServerUserName string
	ServerIP       string
	AppPath        string
	StopCommand    string
	StartCommand   string
	DownloadLogs   bool
	Deploy         bool
}

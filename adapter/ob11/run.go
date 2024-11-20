package ob11

func (config *Config) Run() {
	switch config.Types {
	case "WS":
		config.Driver = NewWebSocketClient(config.URL, config.Token)
	case "WSS":
		//wss有问题
		config.Driver = NewWebSocketServer(16, config.URL, config.Token)
	}
	config.Driver.Connect(config) //连接
	config.Driver.Listen(config)
}

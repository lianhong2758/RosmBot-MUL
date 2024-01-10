package driver

// APICallers 所有的APICaller列表， 通过self-ID映射
var APICallers callerMap

// APICaller is the interface of CallApi
type APICaller interface {
	CallApi(request APIRequest) (APIResponse, error)
}

// Driver 与OneBot通信的驱动，使用driver.DefaultWebSocketDriver
type Driver interface {
	Connect()
	Listen(func([]byte, APICaller))
	GetSelfID() int64
}

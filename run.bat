go version
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
go run main.go
pause
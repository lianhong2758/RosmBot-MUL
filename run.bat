go version
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
go build -ldflags=-checklinkname=0
pause
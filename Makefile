NAME=RosmBot-MUL
EXE_NAME=${NAME}exe
PROTOPATH=server/mys/proto
KANBANPATH=kanban
VERSION=1.1.0

build:
	@echo "build!"
	@go version
	@go env -w GOPROXY=https://goproxy.cn,direct
	@go mod tidy
	@go install github.com/tc-hib/go-winres@latest
	@cd ${KANBANPATH} &&  go-winres make
	@go build -ldflags=-checklinkname=0
	@cd ${KANBANPATH} && rm *.syso
	@echo "Done!"

run:
	@echo "run!"
	@go version
	@go env -w GOPROXY=https://goproxy.cn,direct
	@go mod tidy
	@go run main.go -ldflags=-checklinkname=0

debug:
	@echo "debug"
	@go version
	@go env -w GOPROXY=https://goproxy.cn,direct
	@go mod tidy
	@go run main.go -d -ldflags=-checklinkname=0
	
build_proto:
	@echo "build_proto!"
	@cd ${PROTOPATH} && protoc --go_out=. --go_opt=paths=source_relative *.proto
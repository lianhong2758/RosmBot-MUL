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
	@if ! command -v go-winres &> /dev/null; then \
        echo "go-winres not found. Installing..."; \
       go install github.com/tc-hib/go-winres@latest; \
    fi
	@cd ${KANBANPATH} &&  go-winres make
	@go build
	@cd ${KANBANPATH} && rm *.syso
	@echo "Done!"

run:
	@echo "run!"
	@go version
	@go env -w GOPROXY=https://goproxy.cn,direct
	@go mod tidy
	@go run main.go

debug:
	@echo "debug"
	@go version
	@go env -w GOPROXY=https://goproxy.cn,direct
	@go mod tidy
	@go run main.go -d
	
build_proto:
	@echo "build_proto!"
	@cd ${PROTOPATH} && protoc --go_out=. --go_opt=paths=source_relative *.proto
NAME=RosmBot-MUL
EXE_NAME=${NAME}exe
PROTOPATH=server/mys/proto
VERSION=1.0.0

build:
	@echo "build!"
	@go build

run:
	@echo "run!"
	@go run main.go

debug:
	@echo "debug"
	@go run main.go -d
	
build_proto:
	@echo "build_proto!"
	@cd ${PROTOPATH} && protoc --go_out=. --go_opt=paths=source_relative *.proto
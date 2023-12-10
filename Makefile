export GO111MODULE=on
PROJECT=kp
NAME=be-kredit-pintar
TAG := $(shell git describe --candidates=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
VERSION := $(TAG:v%=%)
ifeq ($(TAG),)
	VERSION := $(BRANCH)
endif

init:
	git config core.hooksPath .githooks
	go mod tidy

build:
	protoc --plugin=protoc-gen-go=${HOME}/go/bin/protoc-gen-go --plugin=protoc-gen-go_grpc=${HOME}/go/bin/protoc-gen-go-grpc --go_out=auth/delivery/grpc --go_grpc_out=auth/delivery/grpc --proto_path=auth/delivery/grpc/proto auth/delivery/grpc/proto/insurance-auth.proto
	go mod tidy
	go build -o ${NAME} app/*.go

clean:
	if [ -f auth/delivery/http/openapi/openapi.yaml ] ; then rm auth/delivery/http/openapi/openapi.yaml ; fi
	if [ -f auth/delivery/http/grpcdoc/insurance-auth.proto ] ; then rm auth/delivery/http/grpcdoc/grpc.proto ; fi
	if [ -f ${NAME} ] ; then rm ${NAME} ; fi

docker:
	docker build -t ${REGISTRY}/${PROJECT}/${NAME}:$(VERSION) .

run:
	protoc --plugin=protoc-gen-go=${HOME}/go/bin/protoc-gen-go --plugin=protoc-gen-go_grpc=${HOME}/go/bin/protoc-gen-go-grpc --go_out=auth/delivery/grpc --go_grpc_out=auth/delivery/grpc --proto_path=auth/delivery/grpc/proto auth/delivery/grpc/proto/insurance-auth.proto
	go mod tidy
	go run app/*.go -c config.yaml

push:
	docker push ${REGISTRY}/${PROJECT}/${NAME}:$(VERSION)

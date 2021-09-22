VERSION = "0.0.1"
change-version:
	@echo $(VERSION)>VERSION
	@echo "package constant\n\n//Version constant of abesh\nconst Version = \"$(VERSION)\"">constant/version.go

update-module:
	go get -v github.com/spf13/cobra
	go get -v go.uber.org/zap
	go get -v github.com/caarlos0/env
	go get -v gopkg.in/yaml.v2
	go get -v google.golang.org/protobuf/proto

protoc:
	@protoc \
		-I=./proto \
			--go_opt=module=github.com/mkawserm/abesh \
			--go_out=. \
			./proto/model/metadata.proto \
			./proto/model/event.proto


run-default:
	go run main/default/main.go run --manifest example/manifest.yaml

build-default:
	go build -o bin/abesh main/default/main.go

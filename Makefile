VERSION = "0.0.1"
change-version:
	@echo $(VERSION)>VERSION
	@echo "package constant\n\n//Version constant of abesh\nconst Version = \"$(VERSION)\"">constant/version.go

update-module:
	go get -v github.com/spf13/cobra
	go get -v go.uber.org/zap
	go get -v github.com/caarlos0/env
	go get -v gopkg.in/yaml.v2

run:
	go run main.go

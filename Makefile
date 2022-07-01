VERSION = "0.0.1"
change-version:
	@echo $(VERSION)>VERSION
	@echo "package constant\n\n//Version constant of abesh\nconst Version = \"$(VERSION)\"">constant/version.go
	@git add VERSION
	@git add constant/version.go
	@git commit -m "v$(VERSION)"
	@git tag -a "v$(VERSION)" -m "v$(VERSION)"
	@git push origin
	@git push origin "v$(VERSION)"

update-module:
	go get -v github.com/spf13/cobra
	go get -v go.uber.org/zap
	go get -v github.com/caarlos0/env/v6
	go get -v gopkg.in/yaml.v2
	go get -v google.golang.org/protobuf/proto

protoc:
	@protoc \
		-I=./proto \
			--go_opt=module=github.com/mkawserm/abesh \
			--go_out=. \
			./proto/model/metadata.proto \
			./proto/model/event.proto \
			./proto/model/status.proto \
			./proto/model/error.proto


build-default:
	go build -o bin/abesh main/default/main.go

run-default:
	go run main/default/main.go run --manifest example/manifest.yaml

run-embedded-print-manifest:
	go run main/embedded/main.go embedded print-manifest

run-embedded:
	go run main/embedded/main.go embedded run

run-embedded2:
	go run main/embedded/main.go embedded-run2

run-race:
	go run -race main/embedded/main.go embedded run

test:
	go test ./... -v

cover:
	go test ./... -coverprofile=cover.out -v

cover-html:
	go tool cover -html=cover.out

profile-heap:
	go tool pprof http://localhost:6060/debug/pprof/heap

profile-cpu:
	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

profile-block:
	go tool pprof http://localhost:6060/debug/pprof/block

profile-mutex:
	go tool pprof http://localhost:6060/debug/pprof/mutex

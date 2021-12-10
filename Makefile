

Workdir ?= internal/demo

demo.up: dockerize
	cd $(Workdir) && go run .

dockerize:
	cd $(Workdir) && go run . dockerize

build: clean
	go build -o out/jarvis cmd/jarvis/main.go

new.project: build
	cd out/ && ./jarvis new --name srv-app --dir somepath

clean:
	rm -rf out

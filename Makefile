
demo.run:
	cd internal/demo && go run .

build: clean
	go build -o out/jarvis cmd/jarvis/main.go

new.project: build
	cd out/ && ./jarvis new --name srv-app

clean:
	rm -rf out
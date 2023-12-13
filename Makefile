dev:
	@go run .

build-all:
	@GOOS=linux GOARCH=amd64 go build -o ./bin/gohooked_amd64
	@GOOS=windows GOARCH=amd64 go build -o ./bin/gohooked_win.exe
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/gohooked_macos_arm
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/gohooked_macos_amd64

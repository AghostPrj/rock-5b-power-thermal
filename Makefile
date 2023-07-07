build : getDep clean arm64

getDep:
	go get -v -t -d ./...


clean :
	rm -rf target; go clean; go clean --cache


arm64 :
	go clean; env GOOS=linux GOARCH=arm64 go build -ldflags "-s -w " -o target/linux-arm64/rock-5b-power-thermal ./cmd/rock-5b-power-thermal


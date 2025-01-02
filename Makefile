
build:
	go build -o ./bin/todor

run:
	./bin/todor


new-release:
	./scripts/version-bump.sh

build-dist:
	OOS=linux GOARCH=amd64 go build -o ./bin/todor-linux-amd64
	OOS=linux GOARCH=386 go build -o ./bin/todor-linux-386i
	OOS=linux GOARCH=arm64 go build -o ./bin/todor-linux-arm64
	OOS=linux GOARCH=arm go build -o ./bin/todor-linux-arm

	OOS=darwin GOARCH=amd64 go build -o ./bin/todor-darwin-amd64.dmg
	OOS=darwin GOARCH=386 go build -o ./bin/todor-darwin-386i.dmg
	OOS=darwin GOARCH=arm64 go build -o ./bin/todor-darwin-arm64.dmg
	OOS=darwin GOARCH=arm go build -o ./bin/todor-darwin-arm.dmg

	OOS=windows GOARCH=amd64 go build -o ./bin/todor-windows-amd64.exe
	OOS=windows GOARCH=386 go build -o ./bin/todor-windows-386i.exe

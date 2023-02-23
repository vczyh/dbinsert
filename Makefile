.PHONY: release, build

release: $(verson)
	@echo "Release Version: $(version)"
	@echo $(version) > VERSION
	git add VERSION
	git commit -m $(version)
	git tag $(version)

beta: $(verson)
	@echo "Beta Version: $(version)-beta"
	@echo $(version)-beta > VERSION
	git add VERSION
	git commit -m $(version)-beta

build:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dbinsert-linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/dbinsert-linux-arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/dbinsert-darwin-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/dbinsert-darwin-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/dbinsert-windows-amd64
	@#CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/dbinsert-windows-arm64

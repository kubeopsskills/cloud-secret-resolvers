TARGET_PATH = bin
GOARCH = GOARCH=amd64
VERSION = 1.2.1
GOMODULE = github.com/kubeopsskills/cloud-secret-resolvers/cmd/csr

buildWindows:
	env GOOS=windows $(GOARCH) go build -o ./$(TARGET_PATH)/csr.exe $(GOMODULE)
	cd $(TARGET_PATH) && zip csr-Windows-$(VERSION).zip ./csr.exe

buildMacOS:
	env GOOS=darwin $(GOARCH) go build  -o ./$(TARGET_PATH)/csr $(GOMODULE)
	cd $(TARGET_PATH) && tar -zcvf csr-MacOS-$(VERSION).tar.gz ./csr

buildLinux:
	env GOOS=linux $(GOARCH) go build -o ./$(TARGET_PATH)/csr $(GOMODULE)
	cd $(TARGET_PATH) && tar -zcvf csr-Linux-amd64-$(VERSION).tar.gz ./csr

buildARM:
	env GOOS=linux GOARCH=arm64 go build -o ./$(TARGET_PATH)/csr $(GOMODULE)
	cd $(TARGET_PATH) && tar -zcvf csr-Linux-arm64-$(VERSION).tar.gz ./csr

build: buildWindows buildMacOS buildLinux buildARM

run: 
	go run cmd/csr/csr.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverpkg=./... -coverprofile=coverage.out ./... && go tool cover -func coverage.out

html-coverage: 
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

clean:
	rm -rf bin

all: clean build

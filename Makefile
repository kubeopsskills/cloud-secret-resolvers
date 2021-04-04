TARGET_PATH = bin
GOARCH = GOARCH=amd64
VERSION = 1.0.0
GOMODULE = github.com/kubeopsskills/cloud-secret-resolvers/cmd/csr

buildWindows:
	env GOOS=windows $(GOARCH) go build -o ./$(TARGET_PATH)/windows/csr.exe $(GOMODULE)
	cd $(TARGET_PATH) && zip csr-Windows-$(VERSION).zip ./windows/csr.exe

buildMacOS:
	env GOOS=darwin $(GOARCH) go build  -o ./$(TARGET_PATH)/macos/csr $(GOMODULE)
	cd $(TARGET_PATH) && tar -zcvf csr-MacOS-$(VERSION).tar.gz ./macos/csr

buildLinux:
	env GOOS=linux $(GOARCH) go build -o ./$(TARGET_PATH)/linux/csr $(GOMODULE)
	cd $(TARGET_PATH) && tar -zcvf csr-Linux-$(VERSION).tar.gz ./linux/csr

buildARM:
	env GOOS=linux GOARCH=arm64 go build -o ./$(TARGET_PATH)/arm/csr $(GOMODULE)
	cd $(TARGET_PATH) && tar -zcvf csr-ARM-$(VERSION).tar.gz ./arm/csr

build: buildWindows buildMacOS buildLinux buildARM

clean:
	rm -rf bin

all: clean build

TARGET_PATH = bin
GOARCH = GOARCH=amd64
VERSION = 1.0.0-alpha

buildWindows:
	env GOOS=windows $(GOARCH) go build -o $(TARGET_PATH)/csr-Windows-$(VERSION).exe
	cd $(TARGET_PATH) && zip csr-Windows-$(VERSION).zip csr-Windows-$(VERSION).exe

buildMacOS:
	env GOOS=darwin $(GOARCH) go build -o $(TARGET_PATH)/csr-MacOS-$(VERSION)
	cd $(TARGET_PATH) && tar -zcvf csr-MacOS-$(VERSION).tar.gz csr-MacOS-$(VERSION)

buildLinux:
	env GOOS=linux $(GOARCH) go build -o $(TARGET_PATH)/csr-Linux-$(VERSION)
	cd $(TARGET_PATH) && tar -zcvf csr-Linux-$(VERSION).tar.gz csr-Linux-$(VERSION)

build: buildWindows buildMacOS buildLinux

clean:
	rm -rf bin

all: clean build
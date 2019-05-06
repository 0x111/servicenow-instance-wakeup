binary_name = servicenow-instance-wakeup

.PHONY: all
all: linux_amd64 darwin_amd64 windows_amd64 checksums

.PHONY: linux_amd64
linux_amd64:
	GOOS=linux GOARCH=amd64 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/$(binary_name)-linux-amd64

.PHONY: linux_i386
linux_i386:
	GOOS=linux GOARCH=386 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/$(binary_name)-linux-i386

.PHONY: darwin_amd64
darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/$(binary_name)-darwin-amd64

.PHONY: darwin_i386
darwin_i386:
	GOOS=darwin GOARCH=386 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/$(binary_name)-darwin-i386

.PHONY: windows_amd64
windows_amd64:
	CC=/usr/local/bin/x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/$(binary_name)-windows-amd64.exe

.PHONY: windows_i386
windows_i386:
	CC=/usr/local/bin/x86_64-w64-mingw32-gcc GOOS=windows GOARCH=386 go build -v -a -gcflags=-trimpath=$$PWD -asmflags=-trimpath=$$PWD -o build/$(binary_name)-windows-i386.exe

.PHONY: checksums
checksums:
	shasum -a 256 build/* > build/checksum.txt

test:
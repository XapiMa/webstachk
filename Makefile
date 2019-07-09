.PHONY: all
all:
	make webStatusChecker_linux
	make webStatusChecker.exe

.PHONY: webStatusChecker_linux
webStatusChecker_linux:
	GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o $@ ./cmd/webStatusChecker
	
.PHONY: webStatusChecker.exe
webStatusChecker.exe:
	GOOS=windows go build -ldflags '-w -s -extldflags "-static"' -o $@ ./cmd/webStatusChecker

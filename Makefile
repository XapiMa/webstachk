all: webStatusChecker webStatusChecker.exe

webStatusChecker:
	GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o $@ ./

webStatusChecker.exe:
	GOOS=windows go build -ldflags '-w -s -extldflags "-static"' -o $@ ./

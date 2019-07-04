.PHONY: flos-hortus
webStatusChecker:
	GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o $@

.PHONY: webStatusChecker.exe
webStatusChecker.exe:
	GOOS=windows go build -ldflags '-w -s -extldflags "-static"' -o $@ 

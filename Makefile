.PHONY: all
all:
	make webstachk_linux
	make webstachk.exe
	make webstachk_mac


.PHONY: webstachk_mac
webstachk_mac:
	GOOS=darwin go build -ldflags '-w -s -extldflags "-static"' -o $@ ./cmd/webstachk

.PHONY: webstachk_linux
webstachk_linux:
	GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o $@ ./cmd/webstachk
	
.PHONY: webstachk.exe
webstachk.exe:
	GOOS=windows go build -ldflags '-w -s -extldflags "-static"' -o $@ ./cmd/webstachk

.PHONY: flos-hortus
flos-hortus:
	GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o $@

.PHONY: flos-hortus.exe
flos-hortus.exe:
	GOOS=windows go build -ldflags '-w -s -extldflags "-static"' -o $@

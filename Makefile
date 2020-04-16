BIN		:= bin
GOCMD	:=go
GOBUILD	:=$(GOCMD) build
GOCLEAN	:=$(GOCMD) clean

.PHONY: all clean
all: proxy

proxy: proxy/proxy.go
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/proxy -v $^

clean:
	$(GOCLEAN)
	rm -f $(BIN)/proxy
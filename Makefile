BIN		:= bin
GOCMD	:=go
GOBUILD	:=$(GOCMD) build
GOCLEAN	:=$(GOCMD) clean

.PHONY: all clean
all: proxy recipient

proxy: proxy/proxy.go
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/proxy -v $^

recipient: recipient/recipient.go
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/recipient -v $^

clean:
	$(GOCLEAN)
	rm -f $(BIN)/proxy
	rm -f $(BIN)/recipient

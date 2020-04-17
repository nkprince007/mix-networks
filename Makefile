BIN		:= bin
GOCMD	:=go
GOBUILD	:=$(GOCMD) build
GOCLEAN	:=$(GOCMD) clean
GOTEST	:=$(GOCMD) test

.PHONY: all clean test proxy recipient
all: proxy recipient

proxy:
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/proxy -v ./proxy

recipient:
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/recipient -v ./recipient

test:
	openssl genrsa -out mixes/test_private.pem 2048
	openssl rsa -in mixes/test_private.pem -outform PEM -pubout -out mixes/test_public.pem
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BIN)/proxy
	rm -f $(BIN)/recipient
	rm -f mixes/test_private.pem mixes/test_public.pem

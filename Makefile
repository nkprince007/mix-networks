BIN		:= bin
GOCMD	:=go
GOBUILD	:=$(GOCMD) build
GOCLEAN	:=$(GOCMD) clean
GOTEST	:=$(GOCMD) test

.PHONY: all clean test
all: proxy recipient

proxy: proxy/proxy.go
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/proxy -v $^

recipient: recipient/recipient.go
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/recipient -v $^

test:
	openssl genrsa -out mixes/test_private.pem 2048
	openssl rsa -in mixes/test_private.pem -outform PEM -pubout -out mixes/test_public.pem
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BIN)/proxy
	rm -f $(BIN)/recipient
	rm -f mixes/test_private.pem mixes/test_public.pem

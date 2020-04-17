BIN		:= bin
GOCMD	:= go
GOBUILD	:= $(GOCMD) build
GOCLEAN	:= $(GOCMD) clean
GOTEST	:= $(GOCMD) test
KEYSIZE := 2048

.PHONY: all clean test proxy recipient sender
all: proxy recipient sender

proxy:
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/proxy -v ./proxy

recipient:
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/recipient -v ./recipient

sender:
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/sender -v ./sender

test:
	openssl genrsa -out mixes/test_private.pem $(KEYSIZE)
	openssl rsa -in mixes/test_private.pem -outform PEM -pubout -out mixes/test_public.pem
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BIN)/proxy
	rm -f $(BIN)/recipient
	rm -f $(BIN)/sender
	rm -f mixes/test_private.pem mixes/test_public.pem

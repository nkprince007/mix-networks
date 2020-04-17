BIN		:= bin
GOCMD	:= go
GOBUILD	:= $(GOCMD) build
GOCLEAN	:= $(GOCMD) clean
GOTEST	:= $(GOCMD) test
KEYSIZE := 2048

.PHONY: all clean test proxy recipient sender
all: proxy recipient sender

proxy: keys
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/proxy -v ./proxy

recipient: keys
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/recipient -v ./recipient

sender: keys
	CGO_ENABLED=0 GOARCH=amd64 $(GOBUILD) -o $(BIN)/sender -v ./sender

test:
	openssl genrsa -out mixes/test_private.pem $(KEYSIZE)
	openssl rsa -in mixes/test_private.pem -outform PEM -pubout -out mixes/test_public.pem
	$(GOTEST) -v ./...

keys:
	openssl genrsa -out recipient/recipient-privkey.pem $(KEYSIZE)
	openssl rsa -in recipient/recipient-privkey.pem -outform PEM -pubout -out sender/recipient-pubkey.pem
	openssl genrsa -out proxy/proxy-privkey.pem $(KEYSIZE)
	openssl rsa -in proxy/proxy-privkey.pem -outform PEM -pubout -out sender/proxy-pubkey.pem

clean:
	$(GOCLEAN)
	rm -f $(BIN)/proxy
	rm -f $(BIN)/recipient
	rm -f $(BIN)/sender
	rm -f mixes/test_private.pem mixes/test_public.pem

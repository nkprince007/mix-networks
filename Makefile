BIN := bin

.PHONY: all
all: proxy

proxy: proxy/proxy.go
	go build -o $(BIN)/proxy -v $^

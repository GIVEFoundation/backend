all: give_api

give_api: main.go
	go build -o give_api main.go

test: 
	go test -v

get_libs:
	go get github.com/ethereum/go-ethereum
	go get github.com/onrik/ethrpc
	go get github.com/gin-gonic/gin
	go get github.com/itsjamie/gin-cors
	go get github.com/sirupsen/logrus
	go get github.com/spf13/viper
	go get github.com/lib/pq
	go get github.com/appleboy/gofight
	go get github.com/buger/jsonparser
	go get github.com/stretchr/testify/assert

clean: 
	rm -f give_api give_api.log upload/*


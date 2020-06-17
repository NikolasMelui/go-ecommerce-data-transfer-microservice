.PHONY: config run build bin

config:
	cp ./cconfig/cconfig.example.go ./cconfig/cconfig.go

run:
	go run ./main.go

build:
	go build

bin:
	./go-ecommerce-data-transfer-microservice

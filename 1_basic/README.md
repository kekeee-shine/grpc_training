# grpc_training

```shell
cd service 
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  -I=. -I={gopath/src}  product.proto

```
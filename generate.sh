#/ktrain:
protoc -I proto/ --go_out=proto/ proto/*.proto --go-grpc_out=proto/
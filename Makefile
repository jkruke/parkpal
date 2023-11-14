.PHONY: protos

protos:
	protoc protos/currency.proto --go_out=. --go-grpc_out=.


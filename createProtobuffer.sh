#!/bin/sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 
cd internal/proto
mkdir pb &> /dev/null
for datei in *.proto; do
  echo "Erstelle Proto Buffer für: $datei"
  protoc $datei --go_out=./pb/ --go-grpc_out=./pb/ --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative 
done
go mod tidy &> /dev/null
protoc --proto_path=./internal/proto --go_opt=paths=source_relative --go_out=plugins=grpc:./internal/pb ./internal/proto/*.proto
cp ./internal/proto/* ./../rubumo.io_frontend/src/proto
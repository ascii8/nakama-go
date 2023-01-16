package proto

//go:generate protoc -I. --go_out=. --go_opt=paths=source_relative proto.proto
//go:generate protoc -I. --go_out=../../ --go_opt=paths=source_relative --go_opt=Mproto.proto=github.com/ascii8/nakama-go;nakama proto.proto
//go:generate mv ../../proto.pb.go ../../proto_test.pb.go

#/!bin/bash

protoc proto/recordpb/record.proto --go_out=plugins=grpc:..
protoc proto/messagepb/message.proto --go_out=plugins=grpc:..
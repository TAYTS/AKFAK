#/!bin/bash

protoc proto/recordpb/record.proto --go_out=plugins=grpc:..
protoc proto/messagepb/message.proto --go_out=plugins=grpc:..
protoc proto/metadatapb/metadata.proto --go_out=plugins=grpc:..
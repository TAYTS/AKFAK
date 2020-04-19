#/!bin/bash

protoc proto/recordpb/record.proto --go_out=plugins=grpc:..
protoc proto/producepb/produce.proto --go_out=plugins=grpc:..
protoc proto/metadatapb/metadata.proto --go_out=plugins=grpc:..
protoc proto/adminpb/admin.proto --go_out=plugins=grpc:..
protoc proto/clientpb/client.proto --go_out=plugins=grpc:..
protoc proto/commonpb/*.proto --go_out=plugins=grpc:..
protoc proto/adminclientpb/*.proto --go_out=plugins=grpc:..
protoc proto/zkmessagepb/*.proto --go_out=plugins=grpc:..
protoc proto/zookeeperpb/*.proto --go_out=plugins=grpc:..
protoc proto/clustermetadatapb/*.proto --go_out=plugins=grpc:..
protoc proto/heartbeatspb/*.proto --go_out=plugins=grpc:..
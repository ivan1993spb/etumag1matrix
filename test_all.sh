#!/bin/bash

go test
go install

addrs=":2020 :3030 :4050 :6010 :8080 :9090"

cd server
go test
go build

for addr in $addrs; do
	echo "start 127.0.0.1$addr"
	./server --http=$addr &
done

cd ../test
go build
echo "RUN TEST..."
./test $addrs

for addr in $addrs; do
	echo "stop 127.0.0.1$addr"
	set `ps -f | grep server.*$addr`
	kill $2
done

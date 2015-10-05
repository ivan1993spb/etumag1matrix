#!/bin/bash

goimports -w=true test/*go server2/*go server/*go *go



go test
go install



cd server
go test
go build

addrs=":2020 :3030 :4050 :6010 :8080 :9090"
for addr in $addrs; do
	echo "start 127.0.0.1$addr"
	./server --http=$addr &
done




cd ../test
go build test.go
echo "RUN TEST 1..."
./test $addrs
cd ..



for addr in $addrs; do
	echo "stop 127.0.0.1$addr"
	set `ps -f | grep server.*$addr`
	kill $2
done





addr=:8080

cd server2
go test
go build
echo "start server"
./server2 -http=$addr &



cd ../test
echo "RUN TEST 2"
go build test2.go
./test2 $addr
cd ..


echo "stop server"
set `ps -f | grep server.*$addr`
kill $2


echo remove trash
git clean -f
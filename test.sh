#!/bin/bash

# server address
addr=:8081

# matrix size
cols=101
rows=101

# count of THREADS
procs=4

# count of GOROUTINES THAT MULTIPLIES
goroutines=6

# # #



goimports -w=true test/*go server2/*go server/*go *go

go test
go install


cd server2
go build
./server2 --http=$addr --procs=$procs --goroutines=$goroutines &

cd ../test
echo "RUN TEST"
go build test2.go
./test2 --addr=$addr --cols=$cols --rows=$rows
cd ..

echo "stop server"
set `ps -f | grep server.*$addr`
kill $2

echo remove trash
git clean -f

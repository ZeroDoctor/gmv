@echo off
cls

cd ..
set GOARCH=amd64
set GOOS=windows

go build -o gmv .

cd build
@echo on

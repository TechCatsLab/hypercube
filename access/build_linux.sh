#!/bin/sh

linuxRelease=hyperaccess.linux

rm -f $linuxRelease

GOOS=linux GOARCH=amd64 go build -o $linuxRelease

scp $linuxRelease root@10.0.0.251:~/hypercube/access
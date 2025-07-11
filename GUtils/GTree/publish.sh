#/bin/sh
go build -ldflags="-s -w"
[ ! -d ~/bin ] && mkdir ~/bin
cp gtree ~/bin

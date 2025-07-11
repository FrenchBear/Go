#/bin/sh
go build -ldflags="-s -w"
[ ! -d ~/bin ] && mkdir ~/bin
cp rwc ~/bin

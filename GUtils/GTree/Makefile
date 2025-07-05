make:	gtree

gtree:	gtree.go reparse.go osspecific_linux.go osspecific_windows.go
		go build $^

clean:
		rm -f gtree

module ggrep

go 1.24.6

require (
	github.com/PieVio/MyGlob v0.0.0-00010101000000-000000000000
	github.com/PieVio/MyMarkup v0.0.0-00010101000000-000000000000
	github.com/PieVio/TextAutoDecode v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.26.0 // indirect
)

replace (
	github.com/PieVio/MyGlob => ../../Packages/MyGlob
	github.com/PieVio/MyMarkup => ../../Packages/MyMarkup
	github.com/PieVio/TextAutoDecode => ../../Packages/TextAutoDecode
)

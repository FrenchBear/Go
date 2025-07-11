module gtree

go 1.24.3

require golang.org/x/text v0.26.0 // direct

require (
	github.com/PieVio/MyMarkup v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.33.0
)

require golang.org/x/term v0.32.0 // indirect

replace github.com/PieVio/MyMarkup => ../../Packages/MyMarkup

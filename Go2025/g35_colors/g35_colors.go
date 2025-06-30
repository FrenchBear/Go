// g35_colors.go
// Learning go, System programming, Colors using fatih/color package
//
// 2025-06-30	PV		First version

package main

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

func main() {
	// Print with default helper functions
	color.Cyan("Prints text in cyan.")
	color.Blue("Prints %s in blue.", "text") // Supports fmt.Printf-like formatting

	// Customize colors and attributes
	red := color.New(color.FgRed).PrintlnFunc()
	red("This is a red message.")

	boldGreen := color.New(color.Bold, color.FgGreen).PrintfFunc()
	boldGreen("This is %s and %s green text.\n", "bold", "bright")

	// RGB colors (if your terminal supports 24-bit colors)
	orange := color.RGB(255, 128, 0).PrintlnFunc()
	orange("This text is orange using RGB.")

	color.RGB(255, 128, 0).Println("foreground orange")
	color.RGB(230, 42, 42).Println("foreground red")

	color.BgRGB(255, 128, 0).Println("background orange")
	color.BgRGB(230, 42, 42).Println("background red")

	// Background colors
	bgBlueFgWhite := color.New(color.BgBlue, color.FgWhite).PrintlnFunc()
	bgBlueFgWhite("White text on blue background.")

	// Sprint functions to get colored strings without printing directly
	yellowString := color.New(color.FgYellow).SprintFunc()
	errorString := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Printf("Here's a %s and an %s.\n", yellowString("warning"), errorString("error"))

	// Disable/Enable colors programmatically
	c := color.New(color.FgCyan)
	c.Println("Prints cyan text")
	c.DisableColor()
	c.Println("This is printed without any color")
	c.EnableColor()
	c.Println("This prints again cyan...")

	// Global disable
	color.NoColor = true // Disables all color output
	color.Green("This won't be green.")
	color.NoColor = false
	color.Green("This will be green again.")

	// Create a custom print function for convenient
	rouge := color.New(color.FgRed).PrintfFunc()
	rouge("warning")
	rouge("error: %s\n", "Some error message")

	// Mix up multiple attributes
	notice := color.New(color.Bold, color.FgGreen).PrintlnFunc()
	notice("don't forget this...")

	/*
	// You can also FprintXxx functions to pass your own io.Writer:
	blue := color.New(color.FgBlue).FprintfFunc()
	blue(myWriter, "important notice: %s", "************************")

	// Mix up with multiple attributes
	success := color.New(color.Bold, color.FgGreen).FprintlnFunc()
	success(myWriter, "don't forget this...")
	*/

	//Or create SprintXxx functions to mix strings with other non-colorized strings:
	yellow := color.New(color.FgYellow).SprintFunc()
	red2 := color.New(color.FgRed).SprintFunc()

	fmt.Printf("this is a %s and this is %s.\n", yellow("warning"), red2("error"))

	info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
	fmt.Printf("this %s rocks!\n", info("package"))
	//Windows support is enabled by default. All Print functions work as intended. However, only for color.SprintXXX
	//functions, user should use fmt.FprintXXX and set the output to color.Output:

	fmt.Fprintf(color.Output, "Windows support: %s\n", color.GreenString("PASS"))

	info2 := color.New(color.FgWhite, color.BgGreen).SprintFunc()
	fmt.Fprintf(color.Output, "this %s rocks!\n", info2("package"))
	//Using with existing code is possible. Just use the Set() method to set the standard output to the given parameters. That way a rewrite of an existing code is not required.

	// Use handy standard colors.
	color.Set(color.FgYellow)

	fmt.Println("Existing text will be now in Yellow")
	fmt.Printf("This one %s\n", "too")

	color.Unset() // don't forget to unset

	// You can mix up parameters
	color.Set(color.FgMagenta, color.Bold)
	defer color.Unset() // use it in your function

	fmt.Println("All text will be now bold magenta.")

	// There might be a case where you want to disable color output (for example to pipe the standard output of your app
	// to somewhere else). `Color` has support to disable colors both globally and for single color definition. For
	// example suppose you have a CLI app and a `--no-color` bool flag. You can easily disable the color output with:
	var flagNoColor = flag.Bool("no-color", false, "Disable color output")

	if *flagNoColor {
		color.NoColor = true // disables colorized output
	}

	// You can also disable the color by setting the NO_COLOR environment variable to any value.
	// It also has support for single color definitions (local). You can disable/enable color output on the fly:
	c = color.New(color.FgCyan)
	c.Println("Prints cyan text")

	c.DisableColor()
	c.Println("This is printed without any color")

	c.EnableColor()
	c.Println("This prints again cyan...")

	c = color.New(color.FgCyan).Add(color.Underline)
	c.Println("Prints cyan text with an underline.")

	c = color.New(color.FgGreen).Add(color.Italic)
	c.Println("Prints green italic text.")

	color.Unset()
	fmt.Println("Back to standard colors")
}

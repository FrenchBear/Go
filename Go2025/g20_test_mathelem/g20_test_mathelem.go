// g20_test_mathelem.go
// Learning go, testing a package
//
// 2025-06-17	PV		First version

package main

import (
	"fmt"

	mathelem "github.com/PieVio/MathElem"
)

func init() {
	fmt.Println("Running init() of main package")
}

func main() {
	fmt.Println("g20_test_mathelem main()")
	fmt.Println("2.5²:", mathelem.Square(2.5))	
	fmt.Println("2.5³:", mathelem.Cube(2.5))	
}

/*
Gemini request:

I'd like to build a learning example of a go package on Windows. I'd like to build a package MathElem below local folder
C:\Development\GitHub\Go\Packages, and in a separate folder below C:\Development\GitHub\Go\Go2025, I'd like to create a
simple test application using MathElem package. To keep things simple, MathElem should only contain two simple
functions:

func Square(x float64) float64 {
    return x*x
}

func Cube(x float64) float64 {
    return x*x*x
}

What are the steps to do that? 
*/

/*

Steps to Create and Use a Go Package
This guide will walk you through setting up a Go package and a separate application that uses it.

Step 1: Create the MathElem Package
First, let's create the directory for your MathElem package and its source file.
Create the directory: C:\Development\GitHub\Go\Packages\MathElem
Inside the MathElem directory, create a file named mathelem.go with the following content:

// Package MathElem provides basic mathematical element operations.
package MathElem

// Square returns the square of a float64 number.
func Square(x float64) float64 {
	return x * x
}

// Cube returns the cube of a float64 number.
func Cube(x float64) float64 {
	return x * x * x
}

Step 2: Initialize the MathElem Module
Navigate to the MathElem directory in your command prompt and initialize it as a Go module. This makes it importable by other projects.
Open your command prompt (e.g., PowerShell or Command Prompt). 
Navigate to the package directory: cd C:\Development\GitHub\Go\Packages\MathElem
Initialize the Go module. You'll need to choose a module path. A common convention is to use your GitHub username and repository name, or just a descriptive path. 
For this example, we'll use github.com/yourusername/MathElem as a placeholder. 
Remember to replace yourusername with your actual GitHub username or a suitable identifier.

go mod init github.com/yourusername/MathElem

(You will see a go.mod file created in this directory.)

Step 3: Create the Test Application
ow, let's create the separate application that will use your MathElem package.
Create the directory: C:\Development\GitHub\Go\Go2025
Inside the Go2025 directory, create a file named main.go with the following content:

package main

import (
	"fmt"
	// Import your MathElem package using the module path you defined
	// Make sure this path matches the 'go mod init' command you used for MathElem.
	// For example: "github.com/yourusername/MathElem"
	mathelem "github.com/yourusername/MathElem"
)

func main() {
	number := 5.0

	// Use the Square function from MathElem
	squaredResult := mathelem.Square(number)
	fmt.Printf("The square of %.2f is %.2f\n", number, squaredResult)

	// Use the Cube function from MathElem
	cubedResult := mathelem.Cube(number)
	fmt.Printf("The cube of %.2f is %.2f\n", number, cubedResult)
}

Important: Make sure the import path github.com/yourusername/MathElem in main.go exactly matches the module path you used when running go mod init in Step 2.

Step 4: Initialize the Test Application Module and Use go mod edit -replaceBecause your MathElem package is local and not yet published, you need to tell your test application where to find it using a replace directive in its go.mod file.
Open your command prompt. Navigate to the application directory: cd C:\Development\GitHub\Go\Go2025
Initialize the Go module for your application: go mod init go2025app (A go.mod file will be created.)
Add a replace directive to tell Go where to find your MathElem package locally. Adjust the paths to reflect your setup.

go mod edit -replace github.com/yourusername/MathElem=C:\Development\GitHub\Go\Packages\MathElem

This command modifies your go.mod file in C:\Development\GitHub\Go\Go2025 to include a line similar to this:

replace github.com/yourusername/MathElem => C:\Development\GitHub\Go\Packages\MathElem

Finally, download dependencies (this will also resolve the local module): go mod tidy

Step 5: Build and Run the Test Application
Now you can build and run your main.go application.
Ensure you are still in the C:\Development\GitHub\Go\Go2025 directory in your command prompt.
Run the application: go run main.go

You should see the following output:
The square of 5.00 is 25.00
The cube of 5.00 is 125.00
*/
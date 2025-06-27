// g30_viper.go
// Learning go, Play with external viper package
// go get github.com/spf13/viper
// https://github.com/spf13/viper
//
// 2025-06-27	PV		First version

package main

import (
	"flag"
	"fmt"
	"os"

	viper "github.com/spf13/viper"
)

// Options
var format string
var verbose bool

func main() {
	flag.StringVar(&format, "f", "", "Specify format of config file, json or yaml")
	flag.BoolVar(&verbose, "v", false, "Verbose option")
	flag.Usage = Usage
	flag.Parse()

	if format == "" {
		Usage()
		os.Exit(1)
	}

	fmt.Println("format:", format)
	fmt.Println("verbose:", verbose)
	fmt.Println()

	switch format {
	case "yaml":
		read_config_yaml()
	case "json":
		read_config_json()
	default:
		fmt.Println("Unsupported format")
	}
}

// Usage overrides default flag version
func Usage() {
	fmt.Println("go_viper, test of viper package")
	flag.PrintDefaults()
}

func read_config_yaml() {
	viper.SetConfigName("viper_config_yaml")       // name of config file (without extension)
	viper.SetConfigType("yaml")                    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(`%localappdata%/BookApps`) // path to look for the config file in
	//viper.AddConfigPath(`C:\Users\Pierr\AppData\Local\BookApps`)   // path to look for the config file in
	viper.AddConfigPath("$HOME/tests_viper") // call multiple times to add many search paths
	viper.AddConfigPath("./config")          // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Println("Yaml config file read successfully")
	fmt.Println("version:", viper.Get("version"))
	fmt.Println("pi:", viper.Get("pi"))
	fmt.Println("command:", viper.Get("command"))
	fmt.Println("description:", viper.Get("description"))

	// Set default value
	viper.SetDefault("speed_limit", 80)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	fmt.Println("speed_limit:", viper.Get("speed_limit"))
	fmt.Println("Taxonomies:", viper.Get("Taxonomies"))

	// Get map
	e := viper.Get("environment")
	if e == nil {
		panic(fmt.Errorf("environment missing"))
	}
	envir, ok := e.(map[string]any)
	if !ok {
		panic(fmt.Errorf("Error casting environment"))
	}
	fmt.Println("envir[source]:", envir["source"])
	fmt.Println("envir[target]:", envir["target"])
	fmt.Println("envir[n1]:", envir["n1"])
	fmt.Println("envir[b1]:", envir["b1"])
	fmt.Println("envir[a1]:", envir["a1"])
	fmt.Println("envir[d1]:", envir["d1"])

}

func read_config_json() {
	viper.SetConfigName("viper_config_json") // name of config file (without extension)
	viper.SetConfigType("json")                    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(`%localappdata%/BookApps`) // path to look for the config file in
	//viper.AddConfigPath(`C:\Users\Pierr\AppData\Local\BookApps`)   // path to look for the config file in
	viper.AddConfigPath("$HOME/tests_viper") // call multiple times to add many search paths
	viper.AddConfigPath("./config")          // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Println("Json config file read successfully")
	fmt.Println("version:", viper.Get("version"))
	fmt.Println("maps:", viper.Get("maps"))
	fmt.Println("masks:", viper.Get("masks"))
	fmt.Println("myDateTime1:", viper.Get("myDateTime1"))
	fmt.Println("fruits:", viper.Get("fruits"))
}

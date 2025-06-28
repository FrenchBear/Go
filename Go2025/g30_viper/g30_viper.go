// g30_viper.go
// Learning go, Play with external viper package
//
// go get github.com/spf13/viper
// go get github.com/go-viper/encoding
// https://github.com/spf13/viper
// https://github.com/go-viper/encoding
//
// https://dev.to/kittipat1413/a-guide-to-configuration-management-in-go-with-viper-5271
//
// 2025-06-27	PV		First version

package main

import (
	"fmt"
	"os"

	pflag "github.com/spf13/pflag"
	viper "github.com/spf13/viper"

	"github.com/go-viper/encoding/hcl"
	"github.com/go-viper/encoding/ini"
	"github.com/go-viper/encoding/javaproperties"
)

func main() {
	v := register_extra_viper_codecs()

	pflag.Usage = Usage

	pflag.BoolP("?", "?", false, "Shows help")
	pflag.BoolP("help", "h", false, "Shows help")
	pflag.StringP("format", "f", "", "Format parameter: yaml, json, toml, ini, hcl or java")
	pflag.BoolP("verbose", "v", false, "Verbose option")
	pflag.CommandLine.SetNormalizeFunc(aliasNormalizeFunc)

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	args := pflag.CommandLine.Args()

	if viper.GetBool("help") || viper.GetBool("?") || len(args) > 0 && (args[0] == "?" || args[0] == "help") {
		Usage()
		os.Exit(0)
	}

	format := viper.GetString("format")
	verbose := viper.GetBool("verbose")

	fmt.Println("format:", format)
	fmt.Println("verbose:", verbose)
	fmt.Println()

	if format == "" {
		Usage()
		os.Exit(1)
	}

	switch format {
	case "yaml":
		read_config_yaml()
	case "json":
		read_config_json()
	case "toml":
		read_config_toml()
	case "ini":
		read_config_ini(v)
	case "hcl":
		read_config_hcl(v)
	case "java":
		read_config_java(v)
	default:
		fmt.Println("Unsupported format")
	}
}

// https://github.com/spf13/viper/blob/master/UPGRADE.md
func register_extra_viper_codecs() *viper.Viper {
	codecRegistry := viper.NewCodecRegistry()

	{
		codec := hcl.Codec{}
		codecRegistry.RegisterCodec("hcl", codec)
		codecRegistry.RegisterCodec("tfvars", codec)
	}

	{
		codec := &javaproperties.Codec{}
		codecRegistry.RegisterCodec("properties", codec)
		codecRegistry.RegisterCodec("props", codec)
		codecRegistry.RegisterCodec("prop", codec)
	}

	codecRegistry.RegisterCodec("ini", ini.Codec{})

	v := viper.NewWithOptions(
		viper.WithCodecRegistry(codecRegistry),
	)

	return v
}

// The aliasNormalizeFunc() function is used for creating additional aliases for flag
// User for long flags, so --fmt is equivalent to --format
func aliasNormalizeFunc(f *pflag.FlagSet, n string) pflag.NormalizedName {
	switch n {
	case "fmt":
		n = "format"
	case "ver":
		n = "verbose"
	}
	return pflag.NormalizedName(n)
}

// Usage overrides default flag version
func Usage() {
	fmt.Println("go_viper, test of viper package")
	//flag.PrintDefaults()
	pflag.PrintDefaults()
}

func read_config_yaml() {
	viper.SetConfigName("viper_config_yaml")       // name of config file (without extension)
	viper.SetConfigType("yaml")                    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(`%localappdata%/BookApps`) // path to look for the config file in
	//viper.AddConfigPath(`C:/Users/Pierr/AppData/Local/BookApps`)   // path to look for the config file in
	viper.AddConfigPath("$HOME/tests_viper") // call multiple times to add many search paths
	viper.AddConfigPath("./config")          // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		panic(fmt.Errorf("fatal error reading yaml config file: %w", err))
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
	viper.SetConfigName("viper_config_json")       // name of config file (without extension)
	viper.SetConfigType("json")                    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(`%localappdata%/BookApps`) // path to look for the config file in
	//viper.AddConfigPath(`C:/Users/Pierr/AppData/Local/BookApps`)   // path to look for the config file in
	viper.AddConfigPath("$HOME/tests_viper") // call multiple times to add many search paths
	viper.AddConfigPath("./config")          // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error reading json config file: %w", err))
	}

	fmt.Println("Json config file read successfully")
	fmt.Println("version:", viper.Get("version"))
	fmt.Println("maps:", viper.Get("maps"))
	fmt.Println("masks:", viper.Get("masks"))
	fmt.Println("myDateTime1:", viper.Get("myDateTime1"))
	fmt.Println("fruits:", viper.Get("fruits"))
}

func read_config_toml() {
	viper.SetConfigName("viper_config_toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(`%localappdata%/BookApps`)
	viper.AddConfigPath("$HOME/tests_viper")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading toml config file: %w", err))
	}

	fmt.Println("Toml config file read successfully")
	fmt.Println("title:", viper.Get("title"))
	fmt.Println("database:", viper.Get("database"))
	fmt.Println("clients.data:", viper.Get("clients.data"))
	fmt.Println("servers.alpha:", viper.Get("servers.alpha"))
}

type IniConfig struct {
	Default_values struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}
	Dev_database struct {
		Port       int    `mapstructure:"port"`
		ForwardX11 bool   `mapstructure:"forwardx11"`
		Name       string `mapstructure:"name"`
	}
}

// fatal error reading ini config file: While parsing config: decoder not found for this format
func read_config_ini(viper *viper.Viper) {
	viper.SetConfigName("viper_config_ini")
	viper.SetConfigType("ini")
	viper.AddConfigPath(`%localappdata%/BookApps`)
	viper.AddConfigPath("$HOME/tests_viper")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading ini config file: %w", err))
	}

	fmt.Println("Ini config file read successfully")
	fmt.Println("default_values:", viper.Get("default_values"))
	fmt.Println("dev_database.name:", viper.Get("dev_database.name"))

	// Unmarshalling
	var config IniConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode into struct, %v", err)
	} else {
		fmt.Printf("IniConfig: %v\n", config)
	}
}

func read_config_hcl(viper *viper.Viper) {
	viper.SetConfigName("viper_config_hcl")
	viper.SetConfigType("hcl")
	viper.AddConfigPath(`%localappdata%/BookApps`)
	viper.AddConfigPath("$HOME/tests_viper")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading hcl config file: %w", err))
	}

	fmt.Println("Hcl config file read successfully")
	fmt.Println("io_mode:", viper.Get("io_mode"))
	fmt.Println("service:", viper.Get("service"))
}

func read_config_java(viper *viper.Viper) {
	viper.SetConfigName("viper_config_java.java")
	viper.SetConfigType("hcl")
	viper.AddConfigPath(`%localappdata%/BookApps`)
	viper.AddConfigPath("$HOME/tests_viper")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading java config file: %w", err))
	}

	fmt.Println("Java config file read successfully: ", viper.ConfigFileUsed())
	fmt.Println("version:", viper.Get("version"))
	fmt.Println("db.username:", viper.Get("db.username"))
	fmt.Println("db.url:", viper.Get("db.url"))
	fmt.Println("db.port:", viper.Get("db.port"))
}

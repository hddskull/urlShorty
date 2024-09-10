package config

import (
	"flag"
	"fmt"
)

// Properties
var LaunchAdr = defaultNetAddress()
var RedirectAdr = defaultNetAddress()

func ConfigureNetAddress() {

	fmt.Println(LaunchAdr.String())
	fmt.Println(RedirectAdr.String())

	flag.Var(LaunchAdr, "a", "Network addres host:port")
	flag.Var(RedirectAdr, "b", "Network addres host:port/shortURL")
	flag.Parse()

	fmt.Println(LaunchAdr.String())
	fmt.Println(RedirectAdr.String())
}

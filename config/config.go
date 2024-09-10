package config

import (
	"flag"
	"fmt"
)

// Properties
var LaunchAdr = defaultNetAddress()
var RedirectAdr = defaultNetAddress()

func ConfigureNetAddress() {
	fmt.Println("before parse")
	fmt.Println(LaunchAdr.String())
	fmt.Println(RedirectAdr.String())

	flag.Var(LaunchAdr, "a", "Network address host:port")
	flag.Var(RedirectAdr, "b", "Network address host:port")
	flag.Parse()

	fmt.Println("after parse")
	fmt.Println(LaunchAdr.String())
	fmt.Println(RedirectAdr.String())
}

package config

import (
	"flag"
)

// Properties
var LaunchAdr = defaultNetAddress()
var RedirectAdr = defaultNetAddress()

func ConfigureNetAddress() {
	flag.Var(LaunchAdr, "a", "Network address host:port")
	flag.Var(RedirectAdr, "b", "Network address host:port")
	flag.Parse()
}

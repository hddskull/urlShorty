package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type netAddress struct {
	Host string
	Port int
}

func defaultNetAddress() *netAddress {
	n := new(netAddress)
	n.Host = "127.0.0.1"
	n.Port = 8080
	return n
}

// check interface implementation
var _ = flag.Value(defaultNetAddress())

func (n *netAddress) String() string {
	return fmt.Sprint(n.Host, ":", n.Port)
}

func (n *netAddress) Set(flagValue string) error {
	vals := strings.Split(flagValue, ":")

	if len(vals) != 2 {
		return errors.New("need address in a form host:port")
	}

	port, err := strconv.Atoi(vals[1])
	if err != nil {
		return err
	}

	n.Host = vals[0]
	n.Port = port

	return nil
}

package server

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Config struct {
	ServerAddress address
}

type address struct {
	host string
	port int
}

func (a *address) String() string {
	return a.host + ":" + strconv.Itoa(a.port)
}

func (a *address) Set(value string) error {
	var err error
	serverAddress := strings.Split(value, ":")
	if len(serverAddress) != 2 {
		return errors.New("need 2 arguments: host and port")
	}
	a.host = serverAddress[0]
	a.port, err = strconv.Atoi(serverAddress[1])

	return err
}

func ParseFlags() {
	addr := address{
		host: "localhost",
		port: 8080,
	}

	flag.Var(&addr, "a", "server address host:port")

	flag.Parse()

	if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
		err := addr.Set(addrEnv)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
	}

	Config.ServerAddress = addr
}

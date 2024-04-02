package config

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

var Config struct {
	ServerAddress  address
	ReportInterval int64
	PollInterval   int64
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
	flag.Int64Var(&Config.ReportInterval, "r", 10, "report interval [sec]")
	flag.Int64Var(&Config.PollInterval, "p", 2, "poll interval [sec]")

	flag.Parse()

	Config.ServerAddress = addr
}

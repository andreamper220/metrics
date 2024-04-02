package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Config struct {
	ServerAddress  address
	ReportInterval int
	PollInterval   int
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
	flag.IntVar(&Config.ReportInterval, "r", 10, "report interval [sec]")
	flag.IntVar(&Config.PollInterval, "p", 2, "poll interval [sec]")

	flag.Parse()

	var err error
	if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
		err = addr.Set(addrEnv)
	}
	if reportIntervalEnv := os.Getenv("REPORT_INTERVAL"); reportIntervalEnv != "" {
		Config.ReportInterval, err = strconv.Atoi(reportIntervalEnv)
	}
	if pollIntervalEnv := os.Getenv("POLL_INTERVAL"); pollIntervalEnv != "" {
		Config.PollInterval, err = strconv.Atoi(pollIntervalEnv)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	Config.ServerAddress = addr
}

package server

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/andreamper220/metrics.git/internal/logger"
)

var Config struct {
	ServerAddress   address
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
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
	flag.IntVar(&Config.StoreInterval, "i", 300, "store to file interval [sec]")
	flag.StringVar(&Config.FileStoragePath, "f", "/tmp/metrics-db.json", "absolute path of file to store")
	flag.BoolVar(&Config.Restore, "r", true, "to restore values from file")
	flag.StringVar(&Config.DatabaseDSN, "d", "", "database DSN")

	flag.Parse()

	var err error
	if addrEnv := os.Getenv("ADDRESS"); addrEnv != "" {
		err = addr.Set(addrEnv)
	}
	if storeIntervalEnv := os.Getenv("STORE_INTERVAL"); storeIntervalEnv != "" {
		Config.StoreInterval, err = strconv.Atoi(storeIntervalEnv)
	}
	if fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH"); fileStoragePathEnv != "" {
		Config.FileStoragePath = fileStoragePathEnv
	}
	if restoreEnv := os.Getenv("RESTORE"); restoreEnv != "" {
		Config.Restore, err = strconv.ParseBool(restoreEnv)
	}
	if databaseDsn := os.Getenv("DATABASE_DSN"); databaseDsn != "" {
		Config.DatabaseDSN = databaseDsn
	}

	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	Config.ServerAddress = addr
}

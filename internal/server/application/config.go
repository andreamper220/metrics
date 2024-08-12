package application

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
	Sha256Key       string
	CryptoKeyPath   string
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

	if flag.Lookup("a") == nil {
		flag.Var(&addr, "a", "server address host:port")
	}
	if flag.Lookup("i") == nil {
		flag.IntVar(&Config.StoreInterval, "i", 300, "store to file interval [sec]")
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&Config.FileStoragePath, "f", "", "absolute path of file to store")
	}
	if flag.Lookup("r") == nil {
		flag.BoolVar(&Config.Restore, "r", true, "to restore values from file")
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&Config.DatabaseDSN, "d", "", "database DSN")
	}
	if flag.Lookup("k") == nil {
		flag.StringVar(&Config.Sha256Key, "k", "", "sha256 key")
	}
	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(&Config.CryptoKeyPath, "crypto-key", "", "path to private key file")
	}

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
	if databaseDsnEnv := os.Getenv("DATABASE_DSN"); databaseDsnEnv != "" {
		Config.DatabaseDSN = databaseDsnEnv
	}
	if sha256KeyEnv := os.Getenv("KEY"); sha256KeyEnv != "" {
		Config.Sha256Key = sha256KeyEnv
	}
	if cryptoKeyPathEnv := os.Getenv("CRYPTO_KEY"); cryptoKeyPathEnv != "" {
		Config.CryptoKeyPath = cryptoKeyPathEnv
	}

	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	Config.ServerAddress = addr
}

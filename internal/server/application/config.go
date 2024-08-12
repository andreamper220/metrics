package application

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andreamper220/metrics.git/internal/logger"
)

type jsonConfig struct {
	Address       string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDSN   string `json:"database_dsn"`
	CryptoKeyPath string `json:"crypto_key"`
}

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
	configFilePath := *flag.String("c", "", "config file path")
	if configFilePathEnv := os.Getenv("CONFIG"); configFilePathEnv != "" {
		configFilePath = configFilePathEnv
	}
	if configFilePath != "" {
		jsonConfigFile, err := os.Open(configFilePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		defer jsonConfigFile.Close()

		byteValue, _ := io.ReadAll(jsonConfigFile)
		var config jsonConfig
		if err = json.Unmarshal(byteValue, &config); err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}

		err = Config.ServerAddress.Set(config.Address)
		Config.Restore = config.Restore
		storeInterval, err := time.ParseDuration(config.StoreInterval)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		Config.StoreInterval = int(storeInterval.Seconds())
		Config.FileStoragePath = config.StoreFile
		Config.DatabaseDSN = config.DatabaseDSN
		Config.CryptoKeyPath = config.CryptoKeyPath
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
	}

	addr := address{
		host: "localhost",
		port: 8080,
	}
	if flag.Lookup("a") == nil {
		flag.Var(&addr, "a", "server address host:port")
		if addr.host != Config.ServerAddress.host || addr.port != Config.ServerAddress.port {
			Config.ServerAddress = addr
		}
	}
	if flag.Lookup("i") == nil {
		var storeInterval int
		flag.IntVar(&storeInterval, "i", 300, "store to file interval [sec]")
		if storeInterval != Config.StoreInterval {
			Config.StoreInterval = storeInterval
		}
	}
	if flag.Lookup("f") == nil {
		var fileStoragePath string
		flag.StringVar(&fileStoragePath, "f", "", "absolute path of file to store")
		if fileStoragePath != Config.FileStoragePath {
			Config.FileStoragePath = fileStoragePath
		}
	}
	if flag.Lookup("r") == nil {
		var restore bool
		flag.BoolVar(&restore, "r", true, "to restore values from file")
		if restore != Config.Restore {
			Config.Restore = restore
		}
	}
	if flag.Lookup("d") == nil {
		var databaseDSN string
		flag.StringVar(&databaseDSN, "d", "", "database DSN")
		if databaseDSN != Config.DatabaseDSN {
			Config.DatabaseDSN = databaseDSN
		}
	}
	if flag.Lookup("k") == nil {
		var sha256Key string
		flag.StringVar(&sha256Key, "k", "", "sha256 key")
		if sha256Key != Config.Sha256Key {
			Config.Sha256Key = sha256Key
		}
	}
	if flag.Lookup("crypto-key") == nil {
		var cryptoKeyPath string
		flag.StringVar(&cryptoKeyPath, "crypto-key", "", "path to private key file")
		if cryptoKeyPath != Config.CryptoKeyPath {
			Config.CryptoKeyPath = cryptoKeyPath
		}
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

package agent

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
)

type jsonConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
	CryptoKeyPath  string `json:"crypto_key"`
}

var Config struct {
	ServerAddress  address
	ReportInterval int
	PollInterval   int
	Sha256Key      string
	RateLimit      int
	CryptoKeyPath  string
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
		reportInterval, err := time.ParseDuration(config.ReportInterval)
		pollInterval, err := time.ParseDuration(config.PollInterval)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		Config.ReportInterval = int(reportInterval.Seconds())
		Config.PollInterval = int(pollInterval.Seconds())
		Config.CryptoKeyPath = config.CryptoKeyPath
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
	if flag.Lookup("r") == nil {
		var reportInterval int
		flag.IntVar(&reportInterval, "r", 10, "report interval [sec]")
		if reportInterval != Config.ReportInterval {
			Config.ReportInterval = reportInterval
		}
	}
	if flag.Lookup("p") == nil {
		var pollInterval int
		flag.IntVar(&pollInterval, "p", 2, "poll interval [sec]")
		if pollInterval != Config.PollInterval {
			Config.PollInterval = pollInterval
		}
	}
	if flag.Lookup("k") == nil {
		var sha256Key string
		flag.StringVar(&sha256Key, "k", "", "sha256 key")
		if sha256Key != Config.Sha256Key {
			Config.Sha256Key = sha256Key
		}
	}
	if flag.Lookup("l") == nil {
		var rateLimit int
		flag.IntVar(&rateLimit, "l", 10, "requests per report")
		if rateLimit != Config.RateLimit {
			Config.RateLimit = rateLimit
		}
	}
	if flag.Lookup("crypto-key") == nil {
		var cryptoKeyPath string
		flag.StringVar(&cryptoKeyPath, "crypto-key", "", "path to public key file")
		if cryptoKeyPath != Config.CryptoKeyPath {
			Config.CryptoKeyPath = cryptoKeyPath
		}
	}

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
	if sha256KeyEnv := os.Getenv("KEY"); sha256KeyEnv != "" {
		Config.Sha256Key = sha256KeyEnv
	}
	if rateLimitEnv := os.Getenv("RATE_LIMIT"); rateLimitEnv != "" {
		Config.RateLimit, err = strconv.Atoi(rateLimitEnv)
	}
	if cryptoKeyPathEnv := os.Getenv("CRYPTO_KEY"); cryptoKeyPathEnv != "" {
		Config.CryptoKeyPath = cryptoKeyPathEnv
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	Config.ServerAddress = addr
}

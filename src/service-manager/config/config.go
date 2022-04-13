package config

import (
	"errors"
	"os"
	"strconv"
)

type ConfigObject struct {
	DBUsername    string
	DBPassword    string
	DBName        string
	DBHost        string
	DBPort        int
	HTTPPort      int
	ParamsPath    string
	SecretsPath   string
	EncryptionKey [32]byte
}

var Config ConfigObject

func LoadConfig() {
	dbUsername := os.Getenv("SERVICE_MANAGER_DB_USERNAME")
	dbPassword := os.Getenv("SERVICE_MANAGER_DB_PASSWORD")
	dbHost := os.Getenv("SERVICE_MANAGER_DB_HOST")
	dbName := os.Getenv("SERVICE_MANAGER_DB_NAME")
	dbPortString := os.Getenv("SERVICE_MANAGER_DB_PORT")
	paramsPath := os.Getenv("SERVICE_MANAGER_PARAMS_PATH")
	secretsPath := os.Getenv("SERVICE_MANAGER_SECRETS_PATH")
	dbPort, err := strconv.Atoi(dbPortString)
	if err != nil {
		panic(err)
	}
	httpPortString := os.Getenv("SERVICE_MANAGER_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		panic(err)
	}
	encryptionKeyString := os.Getenv("SERVICE_MANAGER_ENCRYPTION_KEY")
	if len(encryptionKeyString) < 32 {
		err = errors.New("SERVICE_MANAGER_ENCRYPTION_KEY must be at least 32 characters long")
		panic(err)
	}
	encryptionKey := []byte(encryptionKeyString)
	key := [32]byte{}
	for i := 0; i < 32; i++ {
		key[i] = encryptionKey[i]
	}
	Config.DBUsername = dbUsername
	Config.DBPassword = dbPassword
	Config.DBHost = dbHost
	Config.DBName = dbName
	Config.DBPort = dbPort
	Config.HTTPPort = httpPort
	Config.ParamsPath = paramsPath
	Config.SecretsPath = secretsPath
	Config.EncryptionKey = key
}

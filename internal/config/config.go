package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"os"
	"strconv"
	"sync"
)

type Logger struct {
	Path string
}

type Logs struct {
	Path      string
	PathDb    string
	ReadCycle string
}

type Secret struct {
	AccessSecret string
}
type Server struct {
	Domain string
	MaxAge int
	Port   int
}

type Templates struct {
	Path string
}

type TLS struct {
	Certificate string
	Enable      bool
	Key         string
}

type UsersDB struct {
	PathDb string
}

type Config struct {
	Logger    Logger
	Logs      Logs
	Secret    Secret
	Server    Server
	Templates Templates
	TLS       TLS
	UsersDB   UsersDB
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func GetConfig() (*Config, error) {
	config := &Config{}
	file, err := ioutil.ReadFile("./config/logs.conf")
	if err != nil {
		return &Config{}, erx.New(err)
	}
	err = toml.Unmarshal(file, config)
	if err != nil {
		return &Config{}, erx.New(err)
	}
	return config, nil
}

func SetEnv() error {

	config, err := GetConfig()
	if err != nil {
		return erx.New(err)
	}
	err = os.Setenv("USERS_PATH_DB", (*config).UsersDB.PathDb)
	if err != nil {
		return erx.New(err)
	}
	err = os.Setenv("ACCESS_SECRET", (*config).Secret.AccessSecret)
	if err != nil {
		return erx.New(err)
	}
	err = os.Setenv("LOGGER", (*config).Logger.Path)
	if err != nil {
		return erx.New(err)
	}
	err = os.Setenv("LOGS_PATH", (*config).Logs.Path)
	if err != nil {
		return erx.New(err)
	}
	err = os.Setenv("LOGS_PATH_DB", (*config).Logs.PathDb)
	if err != nil {
		return erx.New(err)
	}
	err = os.Setenv("DOMAIN", (*config).Server.Domain)
	if err != nil {
		return erx.New(err)
	}
	err = os.Setenv("COOKIE_MAX_AGE", fmt.Sprint((*config).Server.MaxAge))
	if err != nil {
		return erx.New(err)
	}

	return err
}

func SetLogfilesEnv(logfiles []models.LogFileStruct) error {

	if logfiles != nil {
		for _, v := range logfiles {
			logfileLength := v.LogfileName + "Length"
			err := os.Setenv(logfileLength, strconv.Itoa(v.FileLength))
			if err != nil {
				return erx.New(err)
			}
			logfileLastSessionDate := v.LogfileName + "LastSessionDate"
			err = os.Setenv(logfileLastSessionDate, strconv.Itoa(int(v.LastSessionDate)))
			if err != nil {
				return erx.New(err)
			}
			logfilePreviousDate := v.LogfileName + "PreviousDate"
			err = os.Setenv(logfilePreviousDate, strconv.Itoa(int(v.PreviousDate)))
			if err != nil {
				return erx.New(err)
			}
			logfileSessionCount := v.LogfileName + "SessionCount"
			err = os.Setenv(logfileSessionCount, strconv.Itoa(v.SessionCount))
			if err != nil {
				return erx.New(err)
			}
		}
		return nil
	}

	return nil
}

func (cfg *Config) RewriteConfig() error {

	file, err := toml.Marshal(cfg)
	if err != nil {
		return erx.New(err)
	}

	err = ioutil.WriteFile("./auth.conf", file, 0644)
	if err != nil {
		return erx.New(err)
	}

	return nil
}

//	PathDb: 		getEnv("PATH_DB", "./data/db1/auth.db"),
//	AccessSecret:	getEnv("ACCESS_SECRET", ""),
//	RefreshSecret:	getEnv("REFRESH_SECRET", ""),
//},
//DebugMode: getEnvAsBool("DEBUG_MODE", true),
//UserRoles: getEnvAsSlice("USER_ROLES", []string{"admin"}, ","),
//MaxUsers:  getEnvAsInt("MAX_USERS", 1),

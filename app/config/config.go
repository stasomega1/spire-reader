package configs

import (
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	ErrParseConfig        = errors.New("error parse config for profile")
	ErrParseConfigEquals0 = errors.New("error parse config for profile, equals 0")
)

type Config struct {
	BindAddr            int    `toml:"bind_addr"`
	LivenessAddr        int    `toml:"liveness_addr"`
	LogLevel            string `toml:"log_level"`
	MaxGoroutinesCount  int    `toml:"max_goroutines_count"`
	JwtKey              []byte
	JwtKeyStr           string `toml:"jwt_key"`
	JwtTokenLiveMinutes int    `toml:"jwt_token_live"`
	PostgresDbConfig    *PostgresDbConfig
}

func NewConfig() *Config {
	return &Config{
		BindAddr:         8080,
		LivenessAddr:     8086,
		LogLevel:         "debug",
		PostgresDbConfig: &PostgresDbConfig{},
	}
}

type PostgresDbConfig struct {
	NeedConnect     bool          `toml:"need_connect"`
	Host            string        `toml:"host"`
	Port            int           `toml:"port"`
	DbName          string        `toml:"db_name"`
	SSLMode         string        `toml:"db_ssl_mode"`
	Username        string        `toml:"db_username"`
	Password        string        `toml:"db_password"`
	MaxIdleConns    int           `toml:"max_idle_conns"`
	MaxOpenConns    int           `toml:"max_open_conns"`
	ConnMaxLifeTime time.Duration `toml:"conn_max_life_time"`
}

func ParseConfig(configPath, profile string, cnf *Config) error {
	log.Printf("Load from profile %s", profile)
	if profile == "docker" {
		var err error
		//Main app part
		cnf.BindAddr, err = ParseIntEnvConfig("BIND_ADDR")
		if err != nil {
			return err
		}
		cnf.LivenessAddr, err = ParseIntEnvConfig("LIVENESS_ADDR")
		if err != nil {
			return err
		}
		cnf.MaxGoroutinesCount, err = ParseIntEnvConfig("PARTNERS_MAX_GOROUTINES_COUNT")
		if err != nil {
			return err
		}
		cnf.LogLevel, err = ParseStrEnvConfig("PARTNERS_LOG_LEVEL")
		if err != nil {
			return err
		}
		cnf.JwtKeyStr, err = ParseStrEnvConfig("JWT_KEY")
		if err != nil {
			return err
		}
		cnf.JwtTokenLiveMinutes, err = ParseIntEnvConfig("JWT_TOKEN_LIVE_MINUTES")
		if err != nil {
			return err
		}
		//DB
		cnf.PostgresDbConfig.Host, err = ParseStrEnvConfig("DB_HOST")
		if err != nil {
			return err
		}
		cnf.PostgresDbConfig.Port, err = ParseIntEnvConfig("DB_PORT")
		if err != nil {
			return err
		}
		cnf.PostgresDbConfig.SSLMode, err = ParseStrEnvConfig("DB_SSLMODE")
		if err != nil {
			return err
		}
		cnf.PostgresDbConfig.DbName, err = ParseStrEnvConfig("DB_NAME")
		if err != nil {
			return err
		}
		cnf.PostgresDbConfig.Username, err = ParseStrEnvConfig("DB_USERNAME_DSR")
		if err != nil {
			return err
		}
		cnf.PostgresDbConfig.Password, err = ParseStrEnvConfig("DB_PASSWORD_DSR")
		if err != nil {
			return err
		}
		cnf.PostgresDbConfig.MaxIdleConns, err = ParseIntEnvConfig("DB_MAX_IDLE_CONN")
		if err != nil {
			return err
		}
		cnf.PostgresDbConfig.MaxOpenConns, err = ParseIntEnvConfig("DB_MAX_OPEN_CONN")
		if err != nil {
			return err
		}
		ConnMaxLifeTimeInt, err := ParseIntEnvConfig("DB_CONN_MAX_LIFE_TIME")
		if err != nil {
			return err
		}
		cnf.PostgresDbConfig.ConnMaxLifeTime = time.Duration(ConnMaxLifeTimeInt) * time.Minute

		tz, err := ParseStrEnvConfig("TZ")
		local, err := time.LoadLocation(tz)
		if err != nil {
			log.Printf("error load tz, env value is %s, err is: %v", tz, err)
			return ErrParseConfig
		}
		time.Local = local
	} else if profile == "local" {
		if _, err := toml.DecodeFile(configPath, cnf); err != nil {
			return err
		}
	}
	return nil
}

func ConfigureDB(cnf *PostgresDbConfig) (*sqlx.DB, error) {

	databaseUrl := cnf.GetConnectionString()
	db, err := sqlx.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}
	if !cnf.NeedConnect {
		return db, nil
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cnf.MaxIdleConns)
	db.SetMaxOpenConns(cnf.MaxOpenConns)
	return db, nil
}

func ParseIntEnvConfig(envName string) (int, error) {
	strConf := os.Getenv(envName)
	if strConf == "" {
		log.Printf("err when parse env variable, varaible is empty, env parameter: %s", envName)
		return 0, ErrParseConfig
	}
	intConf, err := strconv.Atoi(strConf)
	if err != nil {
		log.Printf("err when parse int env variable, err: %v, env parameter: %s", err, envName)
		return 0, ErrParseConfig
	}
	if intConf < 0 {
		log.Printf("err when parse env variable, int parameter less than 0, env parameter: %s", envName)
		return 0, ErrParseConfig
	}
	if intConf == 0 {
		log.Printf("err when parse env variable, int parameter equal 0, env parameter: %s", envName)
		return 0, ErrParseConfigEquals0
	}
	return intConf, nil
}

func ParseStrEnvConfig(envName string) (string, error) {
	strConf := os.Getenv(envName)
	if strConf == "" {
		log.Printf("err when parse env variable, varaible is empty, env parameter: %s", envName)
		return "", ErrParseConfig
	}
	return strConf, nil
}

func ParseBoolEnvConfig(envName string) (bool, error) {
	strConf := os.Getenv(envName)
	if strConf == "" {
		log.Printf("err when parse env variable, varaible is empty, env parameter: %s", envName)
		return false, ErrParseConfig
	}
	bolConf, err := strconv.ParseBool(strConf)
	if err != nil {
		log.Printf("err when parse bool env variable, err: %v, env parameter: %s", err, envName)
		return false, err
	}
	return bolConf, nil
}

func (p *PostgresDbConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s password=%s",
		p.Host,
		p.Port,
		p.DbName,
		p.SSLMode,
		p.Username,
		p.Password)
}

func ConfigureLogger(level string, dateFormat string) (*logrus.Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetLevel(lvl)
	if dateFormat != "" {
		customFormatter := &logrus.TextFormatter{}
		customFormatter.TimestampFormat = dateFormat //"02.01.2006 15:04:05"//"dd.mm.yyyy HH24:MI:SS"
		customFormatter.FullTimestamp = true
		logger.SetFormatter(customFormatter)
	}
	return logger, nil
}

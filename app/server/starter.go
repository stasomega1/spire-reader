package server

import (
	"github.com/sirupsen/logrus"
	configs "spire-reader/app/config"
	"spire-reader/app/services"
	"spire-reader/app/store"
	"sync"
)

const (
	LoggerTimeFormat = "02.01.2006 15:04:05"
)

func Start(config *configs.Config) error {
	/****************************************Configure app**************************/
	logger, err := configureLogger(config.LogLevel)
	if err != nil {
		return err
	}
	db, err := configs.ConfigureDB(config.PostgresDbConfig)
	if err != nil {
		return err
	}
	myStore := store.NewStore(db)
	spireService := services.NewSpireService(myStore, logger)
	//Configure JWT
	config.JwtKey = []byte("")
	if config.JwtKeyStr != "" {
		config.JwtKey = []byte(config.JwtKeyStr)
	}
	//
	srv := NewServer(spireService, logger, config.JwtKey, config.JwtTokenLiveMinutes)
	/***********************************start app**************************************/
	var wg sync.WaitGroup
	errors := make(chan error, 2)
	//StartApiServer
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Infof("Start server on %d port", config.BindAddr)
		err := srv.Start(config.BindAddr)
		//err := http.ListenAndServe(fmt.Sprintf(":%d", config.BindAddr), srv)
		if err != nil {
			logger.Errorf("Err when start api server, err: %v", err)
			errors <- err
		}
	}()
	/********************************************************************************/
	select {
	case err := <-errors:
		return err
	default:
		wg.Wait()
	}
	return nil
}

func configureLogger(level string) (*logrus.Logger, error) {
	return configs.ConfigureLogger(level, LoggerTimeFormat)
}

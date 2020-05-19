package services

import (
	"github.com/sirupsen/logrus"
	"spire-reader/app/store"
)

type SpireService struct {
	store  *store.Store
	logger *logrus.Logger
}

func NewSpireService(store *store.Store, logger *logrus.Logger) *SpireService {
	return &SpireService{
		store:  store,
		logger: logger,
	}
}

func (spireService *SpireService) Version() (string, error) {
	return spireService.store.PgdbRepository().Version()
}

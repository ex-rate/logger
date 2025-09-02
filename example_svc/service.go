package examplesvc

import (
	"errors"
	"exrate/logger/pkg"
)

// Service - сервис для примера
type Service struct {
	logger *pkg.Logger
}

func NewService(logger *pkg.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Run() {
	s.logger.Info("Service started")
	s.logger.Debug("Debug message from service")
	s.logger.Warn("Warning message from service")
	s.logger.Error("Error message from service")

	s.logger.WithField("error", errors.New("test error")).Error("Error message from service")
}

func (s *Service) Shutdown() {
	s.logger.Info("Service shutdown")
}

package core

import (
	"github.com/sirupsen/logrus"
)

type (
	// LoggerInterface provides common log methods. This is replica for logrus.FieldLogger.
	LoggerInterface interface {
		logrus.FieldLogger
	}

	// Loggable determine possibility to inject LoggerInterface. Can be used to wire Source and Registry implementations
	Loggable interface {
		WithLogger(LoggerInterface)
	}
)

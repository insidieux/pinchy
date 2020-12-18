package core

type (
	// LoggerInterface provides common log methods. This is replica for logrus.FieldLogger.
	LoggerInterface interface {
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Printf(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Warningf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Fatalf(format string, args ...interface{})
		Panicf(format string, args ...interface{})

		Debug(args ...interface{})
		Info(args ...interface{})
		Print(args ...interface{})
		Warn(args ...interface{})
		Warning(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Panic(args ...interface{})

		Debugln(args ...interface{})
		Infoln(args ...interface{})
		Println(args ...interface{})
		Warnln(args ...interface{})
		Warningln(args ...interface{})
		Errorln(args ...interface{})
		Fatalln(args ...interface{})
		Panicln(args ...interface{})
	}

	// Loggable determine possibility to inject LoggerInterface. Can be used to wire Source and Registry implementations
	Loggable interface {
		WithLogger(LoggerInterface)
	}
)

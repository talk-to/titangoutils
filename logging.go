package titangoutils

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig represents the configuration for logging.
type LogConfig struct {
	ProjectName         string        `json:"project_name" yaml:"project_name"`
	LogBaseDirectory    string        `json:"log_base_directory" yaml:"log_base_directory"`
	LogRotationInterval time.Duration `json:"log_rotation_interval" yaml:"log_rotation_interval"`
	LogMaxAgeDays       int           `json:"log_max_age_days" yaml:"log_max_age_days"`
	LogCompress         bool          `json:"log_compress" yaml:"log_compress"`
	LogToStdOut         bool          `json:"log_to_std_out" yaml:"log_to_std_out"`
	DebugMode           bool          `json:"debug_mode" yaml:"debug_mode"`
}

// logrus logger object

// const (
// 	projectName         = "test"
// 	logBaseDirectory    = "/logs"
// 	logRotationInterval = 1 * time.Hour
// 	logMaxAgeDays       = 7
// 	logCompress         = true
// )

func initLogger(config *LogConfig) *logrus.Logger {

	// default log level
	var log = logrus.New()

	log.Level = logrus.InfoLevel

	// set log level to debug level if debug is true
	if config.DebugMode {
		log.Level = logrus.DebugLevel
	}

	// set log output to stdout if stdout is true
	if config.LogToStdOut {
		log.Out = os.Stdout
		return log
	}

	logFile := config.ProjectName + ".log"
	logDirectory := config.LogBaseDirectory + "/" + config.ProjectName
	logPath := logDirectory + "/" + logFile

	// create log directory
	if err := os.MkdirAll(logDirectory, 0744); err != nil {
		log.Fatalf("can't create log directory=%s, error=%s", logDirectory, err)
	}

	// create log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to set logging to file=%s, error=%s", logPath, err)
	}

	// log on stdout
	log.WithFields(logrus.Fields{"filename": logPath, "max_age": config.LogMaxAgeDays, "compress": config.LogCompress}).Info("logger settings")

	// set logging to file
	log.Out = file

	// set the lumberjack logger
	lumberjackLogger := &lumberjack.Logger{
		Filename: logPath,
		MaxAge:   config.LogMaxAgeDays,
		Compress: config.LogCompress,
		MaxSize:  1024,
	}

	// Set the Lumberjack logger
	log.SetOutput(lumberjackLogger)

	// disable color coding while logging to file
	log.SetFormatter(&logrus.TextFormatter{ForceColors: false, FullTimestamp: true, TimestampFormat: time.RFC3339Nano})

	if config.LogRotationInterval > 0 {
		go rotateLogsPeriodically(lumberjackLogger, config.LogRotationInterval)
	}
	return log
}

func rotateLogsPeriodically(logger *lumberjack.Logger, interval time.Duration) {
	now := time.Now()
	nextHour := now.Truncate(time.Hour).Add(time.Hour)
	durationUntilNextHour := time.Until(nextHour)

	timer := time.NewTimer(durationUntilNextHour)
	<-timer.C

	rotateLogs(logger)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		rotateLogs(logger)
	}
}

func rotateLogs(logger *lumberjack.Logger) {
	err := logger.Rotate()
	var msg string
	if err != nil {
		msg = fmt.Sprintf("log rotation failed: %v", err)
		fmt.Println(msg)
	} else {
		msg = "log rotated successfully"
	}
	fmt.Println(msg)
}

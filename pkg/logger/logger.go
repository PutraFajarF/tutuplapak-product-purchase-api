package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	cfg *config.Config
}

type Log struct {
	Event        string
	StatusCode   interface{}
	ResponseTime time.Duration
	Method       string
	Response     interface{} `json:"response"`
	Key          string      `json:"key,omitempty"`
	Query        string      `json:"query,omitempty"`
	Request      interface{} `json:"request"`
	Message      interface{} `json:"message"`
}

var logger = logrus.New()

func New(cfg *config.Config) *Logger {
	return &Logger{cfg: cfg}
}

var (
	LVL_ERROR = "error"
	LVL_INFO  = "info"
	LVL_WARN  = "warning"
)

func (l *Logger) CreateLog(data *Log, types string) error {
	logName := fmt.Sprintf("%s/logs/service.log", l.cfg.Log.Path)

	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	logger.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      false,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableQuote:    true,
	})

	if err == nil {
		logger.Out = file
	} else {
		log.Println("err", err)
		logger.Info("Failed to log to file, using default stderr")
	}

	if types == LVL_WARN {
		logger.WithFields(logrus.Fields{
			"event":         data.Event,
			"status_code":   data.StatusCode,
			"response_time": data.ResponseTime,
			"method":        data.Method,
			"response":      data.Response,
			"key":           data.Key,
			"query":         data.Query,
			"request":       data.Request,
		}).Warn(data.Message)
	}

	if types == LVL_INFO {
		logger.WithFields(logrus.Fields{
			"event":         data.Event,
			"status_code":   data.StatusCode,
			"response_time": data.ResponseTime,
			"method":        data.Method,
			"response":      data.Response,
			"key":           data.Key,
			"query":         data.Query,
			"request":       data.Request,
		}).Info(data.Message)
	}

	if types == LVL_ERROR {
		logger.WithFields(logrus.Fields{
			"event":         data.Event,
			"status_code":   data.StatusCode,
			"response_time": data.ResponseTime,
			"method":        data.Method,
			"response":      data.Response,
			"key":           data.Key,
			"query":         data.Query,
			"request":       data.Request,
		}).Error(data.Message)
	}

	logger.Out = os.Stdout

	return nil
}

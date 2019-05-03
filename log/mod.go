package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gec2/opts"
	"os"
	"time"
)

// FileLogHook Hook to send logs via a File
type FileLogHook struct {
	Writer *os.File
}

func NewFileLogHook() (*FileLogHook, error) {
	t := time.Now()
	name := fmt.Sprintf(
		"%s/%d_%02d_%02dT%02d:%02d:%02d_deployment.log",
		opts.Opts.DeployContext,
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)
	return &FileLogHook{file}, err
}

func (hook *FileLogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	_, err = hook.Writer.Write([]byte(line))
	return err
}

func (hook *FileLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var log = logrus.New()
var Info = log.Info
var Infof = log.Infof
var Error = log.Error
var Errorf = log.Errorf
var Fatal = log.Fatal
var Fatalf = log.Fatalf
var SetLevel = log.SetLevel
var WithFields = log.WithFields
var Debugf = log.Debugf

func Setup() error {
	log.Out = os.Stdout
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	hook, err := NewFileLogHook()

	if err != nil {
		log.Error("Unable to open file for logging.")
		return fmt.Errorf("Unable to open file for logging.")
	} else {
		log.AddHook(hook)
		return nil
	}
}

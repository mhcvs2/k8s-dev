package common

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
)

const (
	NOCOLOR = 0
	RED     = 31
	GREEN   = 32
	YELLOW  = 33
	BLUE    = 36
	GRAY    = 37
	
	SHELLCOMMANDKEY = "cmd"
)

type ShellLogFormatter struct {
}

func (f *ShellLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	var levelColor int
	switch entry.Level {
	case logrus.ErrorLevel:
		levelColor = RED
	default:
		levelColor = BLUE
	}

	fmt.Fprintf(b, "\x1b[%dm%s:\x1b[%dm %s\n", GRAY, entry.Data[SHELLCOMMANDKEY], levelColor, entry.Message)
	return b.Bytes(), nil
}
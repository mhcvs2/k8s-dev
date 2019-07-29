package common

import (
	"github.com/Sirupsen/logrus"
	"os/exec"
)

var shellLogger = logrus.New()

func init() {
	shellLogger.SetFormatter(&ShellLogFormatter{})
}

func RunShell(command string) error {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmdLog := *shellLogger.WithFields(logrus.Fields{SHELLCOMMANDKEY: command})
	cmd.Stdout = cmdLog.WriterLevel(logrus.InfoLevel)
	cmd.Stderr = cmdLog.WriterLevel(logrus.ErrorLevel)
	err := cmd.Start()
	if err != nil {
		return err
	}
	if err = cmd.Wait(); err != nil {
		return err
	}
	return  nil
}

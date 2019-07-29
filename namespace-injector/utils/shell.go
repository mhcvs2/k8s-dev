package utils

import (
	"github.com/Sirupsen/logrus"
	"k8s-dev/namespace-injector/common"
	"os/exec"
)

var shellLogger = logrus.New()

func init() {
	shellLogger.SetFormatter(&common.ShellLogFormatter{})
}

func RunShell(command string) error {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmdLog := *shellLogger.WithFields(logrus.Fields{common.SHELLCOMMANDKEY: command})
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

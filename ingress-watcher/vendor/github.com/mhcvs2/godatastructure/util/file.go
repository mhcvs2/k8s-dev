package util

import (
	"bufio"
	"bytes"
	"os"
	"io"
)

func Readlines(filePath string, handler func(line []byte) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for lino := 1; ; lino++ {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		line = bytes.TrimRight(line, "\n\r")
		err = handler(line)
		if err != nil {
			return err
		}
	}
	return nil
}
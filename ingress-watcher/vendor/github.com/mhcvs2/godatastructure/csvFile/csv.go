package csvFile

import (
	"github.com/mhcvs2/godatastructure/util"
	"os"
	"encoding/csv"
	"sync"
)

type CSV struct {
	path string
	initData []string
	writer *csv.Writer
	rwLock sync.RWMutex
	f *os.File
}

func NewCSVFile(path string) *CSV {
	return &CSV{path:path}
}

func (c *CSV) Init(initData ...string) error {
	if exist, err := util.Exists(c.path); err != nil {
		return err
	} else if exist {
		if c.writer != nil {
			return nil
		}
		if f, err := os.OpenFile("/root/github/go/src/csv/test.csv", os.O_WRONLY, 0644); err !=nil {
			return err
		} else {
			c.f = f
			c.writer = csv.NewWriter(f)
		}
		return nil
	} else if f, err := os.Create(c.path); err != nil {
		return err
	} else {
		c.rwLock.Lock()
		defer c.rwLock.Unlock()
		f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
		w := csv.NewWriter(f)//创建一个新的写入文件流
		if len(initData) > 0 {
			w.Write(initData)
			c.initData = initData
		} else if len(c.initData) > 0 {
			w.Write(c.initData)
		}
		c.writer = w
		c.f = f
	}
	return nil
}

func (c *CSV) Write(data ...string) error {
	if err := c.Init(); err !=nil {
		return err
	}
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	c.writer.Write(data)
	return nil
}

func (c *CSV) WriteAll(data [][]string) error {
	if err := c.Init(); err !=nil {
		return err
	}
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	c.writer.WriteAll(data)
	return nil
}

func (c *CSV) Done() {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	c.writer.Flush()
	c.f.Close()
}

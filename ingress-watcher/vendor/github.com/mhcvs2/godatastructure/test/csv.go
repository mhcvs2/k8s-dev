package main

import "github.com/mhcvs2/godatastructure/csvFile"

func main() {
	c := csvFile.NewCSVFile("/root/GoglandProjects/beegoTest/src/github.com/mhcvs2/godatastructure/test.csv")
	c.Init("adsfasdf", "sadfsgdfsg")
	c.Write("hahaA", "lala")
	c.Write("hahaA", "lala")
	c.Write("hahaA", "lala")
	c.Done()
}

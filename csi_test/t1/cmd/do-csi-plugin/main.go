package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"k8s-dev/csi_test/t1/driver"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "unix:///var/lib/kubelet/plugins/"+driver.DriverName+"/csi.sock", "CSI endpoint")
		token    = flag.String("token", "", "DigitalOcean access token")
		url      = flag.String("url", "https://api.digitalocean.com/", "DigitalOcean API URL")
		doTag    = flag.String("do-tag", "", "Tag DigitalOcean volumes on Create/Attach")
		version  = flag.Bool("version", false, "Print the version and exit.")
	)
	flag.Parse()

	if *version {
		fmt.Printf("%s - %s (%s)\n", driver.GetVersion(), driver.GetCommit(), driver.GetTreeState())
		os.Exit(0)
	}

	drv, err := driver.NewDriver(*endpoint, *token, *url, *doTag)
	if err != nil {
		log.Fatalln(err)
	}

	if err := drv.Run(); err != nil {
		log.Fatalln(err)
	}
}

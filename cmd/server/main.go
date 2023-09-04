package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PrasadG193/sample-csi-cbt-service/pkg/driver"
)

func main() {
	var endpoint string
	flag.StringVar(&endpoint, "endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	//flag.StringVar(&cfg.DriverName, "drivername", "hostpath.csi.k8s.io", "name of the driver")
	flag.Parse()

	driver := driver.NewSampleDriver(endpoint)
	if err := driver.Run(); err != nil {
		fmt.Printf("Failed to run driver: %s", err.Error())
		os.Exit(1)
	}
}

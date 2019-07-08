package main

import (
	"fmt"
	"os"

	"github.com/cristim/ec2-instances-info"
)

func main() {

	data, err := ec2instancesinfo.Data()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, i := range *data {
		fmt.Print(
			"Instance type: ", i.InstanceType,
			",\tCPU: ", i.PhysicalProcessor,
			",\tCPU cores: ", i.VCPU,
			",\tMemory(GB): ", i.Memory,
			",\tEBS Throughput(MB/s): ", i.EBSThroughput,
			",\tcost in us-east-1: ", i.Pricing["us-east-1"].Linux.OnDemand,
			",\tcost in eu-central-1: ")

		p := i.Pricing["eu-central-1"].Linux.OnDemand

		if p == 0 {
			fmt.Print("UNAVAILABLE")
		} else {
			fmt.Print(p)
		}

		if i.Storage != nil {
			fmt.Print(",\tLocal storage volume size(GB): ", i.Storage.Size,
				",\tLocal storage volumes: ", i.Storage.Devices,
				",\tLocal storage SSD: ", i.Storage.SSD)
		}

		fmt.Println(",\tEBS surcharge: ", i.Pricing["us-east-1"].EBSSurcharge)
	}

}

package main

import (
	"fmt"
	"os"
	"time"

	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info"
)

func main() {

	key := "API_KEY" // keys available from contact@leanercloud.com

	go ec2instancesinfo.Updater(2, nil, &key) // use 0 or negative values for weekly updates

	time.Sleep(30 * time.Second)

	data, err := ec2instancesinfo.Data()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, i := range *data {
		fmt.Print(
			"Instance type: ", i.InstanceType,
			",\tCPU: ", i.PhysicalProcessor,
			",\t Arch: ", i.Arch[0],
			",\tCPU cores: ", i.VCPU,
			",\tMemory(GB): ", i.Memory,
			",\tEBS Throughput(MB/s): ", i.EBSThroughput,
			",\tLinux OD cost in us-east-1: ", i.Pricing["us-east-1"].Linux.OnDemand,
			",\tWindows OD cost in us-east-1: ", i.Pricing["us-east-1"].MSWin.OnDemand,
			",\tRHEL OD cost in us-east-1: ", i.Pricing["us-east-1"].RHEL.OnDemand,
			",\tSLES OD cost in us-east-1: ", i.Pricing["us-east-1"].SLES.OnDemand,
			",\tLinux Spot cost in us-east-1: ", i.Pricing["us-east-1"].Linux.SpotMin,
			",\tLinux Standard RI 1y AllUpfront cost in us-east-1: ", i.Pricing["us-east-1"].Linux.Reserved.StandardAllUpfront1Year)
		if i.Storage != nil {
			fmt.Print(",\tLocal storage volume size(GB): ", i.Storage.Size,
				",\tLocal storage volumes: ", i.Storage.Devices,
				",\tLocal storage SSD: ", i.Storage.SSD)
		}

		fmt.Println(",\tEBS surcharge: ", i.Pricing["us-east-1"].EBSSurcharge)
	}

}

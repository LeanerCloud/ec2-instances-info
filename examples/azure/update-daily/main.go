package main

import (
	"fmt"
	"os"
	"time"

	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info"
)

func main() {
	key := "API_KEY" // keys available from contact@leanercloud.com

	// Start the Azure updater in a goroutine, it will update data every 2 days
	go ec2instancesinfo.AzureUpdater(2, nil, &key) // use 0 or negative values for weekly updates

	// Give some time for the initial data to be fetched
	time.Sleep(30 * time.Second)

	// Get the updated Azure VM data
	data, err := ec2instancesinfo.AzureData()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	region := "us-east"

	fmt.Println("Azure VM Instance Data (with continuous updates):")
	fmt.Println("------------------------------------------------")

	for _, i := range *data {
		// Only display instances with pricing in the selected region
		if pricing, found := i.Pricing[region]; found {
			fmt.Printf("Instance type: %s\n", i.InstanceType)
			fmt.Printf("  Pretty name: %s (%s)\n", i.PrettyName, i.PrettyNameAzure)
			fmt.Printf("  vCPU: %d\n", i.VCPU)
			fmt.Printf("  Memory: %.2f GB\n", i.Memory)
			fmt.Printf("  Category: %s\n", i.Category)
			fmt.Printf("  Linux On-Demand cost: $%.4f\n", pricing.Linux.OnDemand)
			fmt.Printf("  Windows On-Demand cost: $%.4f\n", pricing.Windows.OnDemand)
			fmt.Printf("  Linux Spot Min cost: $%.4f\n", pricing.Linux.SpotMin)
			fmt.Printf("  Windows Spot Min cost: $%.4f\n", pricing.Windows.SpotMin)

			if i.Storage != nil {
				fmt.Printf("  Storage size: %d bytes, NVMe SSD: %t\n", i.Storage.Size, i.Storage.NVMeSSD)
			}
			fmt.Println("------------------------------------------------")
		}
	}
}

package main

import (
	"fmt"
	"os"

	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info"
)

func main() {
	key := "API_KEY" // keys available from contact@leanercloud.com

	// One-time update of Azure VM data
	err := ec2instancesinfo.UpdateAzureData(nil, &key)
	if err != nil {
		fmt.Println("Couldn't update Azure VM data, reverting to static compile-time data:", err.Error())
	}

	// Get the Azure VM data
	data, err := ec2instancesinfo.AzureData()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Display information about all Azure VM instances in a specific region
	region := "us-east"

	fmt.Println("Azure VM Instance Information (one-off update):")
	fmt.Println("---------------------------------------------")

	for _, i := range *data {
		// Only display the first 10 instances to avoid too much output
		if pricing, found := i.Pricing[region]; found {
			fmt.Printf("Instance type: %s (%s)\n", i.InstanceType, i.PrettyNameAzure)
			fmt.Printf("  Category: %s\n", i.Category)
			fmt.Printf("  vCPU: %d\n", i.VCPU)
			fmt.Printf("  Memory: %.2f GB\n", i.Memory)
			fmt.Printf("  ACU: %d\n", i.ACU)

			// Print pricing for Linux and Windows
			fmt.Printf("  Linux On-Demand cost: $%.4f\n", pricing.Linux.OnDemand)
			fmt.Printf("  Windows On-Demand cost: $%.4f\n", pricing.Windows.OnDemand)

			// Print storage information if available
			if i.Storage != nil {
				fmt.Printf("  Storage: %d bytes\n", i.Storage.Size)
			}

			fmt.Println("---------------------------------------------")
		}
	}
}

package main

import (
	"fmt"
	"os"
	"strings"

	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info"
)

func main() {
	data, err := ec2instancesinfo.AzureData()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	instanceFamily := "a" // Use the 'a' family as shown in the JSON example
	region := "us-east"   // Adjust region as needed

	fmt.Println("Azure VM Instance Information:")
	fmt.Println("-----------------------------")

	for _, i := range *data {
		// Filter for the specified instance family
		if !strings.HasPrefix(i.InstanceType, instanceFamily) {
			continue
		}

		fmt.Printf("\nInstance Type: %s (%s)\n", i.InstanceType, i.PrettyNameAzure)
		fmt.Printf("Category: %s\n", i.Category)
		fmt.Printf("Family: %s\n", i.Family)
		fmt.Printf("Architecture: %s\n", strings.Join(i.Arch, ", "))
		fmt.Printf("vCPUs: %d (Available: %d, Per Core: %d)\n", i.VCPU, i.VCPUsAvailable, i.VCPUsPerCore)
		fmt.Printf("Memory: %.2f GB\n", i.Memory)
		fmt.Printf("ACU: %d\n", i.ACU)
		fmt.Printf("Accelerated Networking: %t\n", i.AcceleratedNetworking)
		fmt.Printf("Premium IO: %t\n", i.PremiumIO)

		// Check if the region is available in pricing
		if pricing, found := i.Pricing[region]; found {
			fmt.Printf("\nPricing in %s:\n", region)
			fmt.Printf("  Linux:\n")
			fmt.Printf("    On-Demand: $%.4f\n", pricing.Linux.OnDemand)
			fmt.Printf("    Basic: $%.4f\n", pricing.Linux.Basic)
			fmt.Printf("    Basic Spot: $%.4f\n", pricing.Linux.BasicSpot)
			fmt.Printf("    Spot Min: $%.4f\n", pricing.Linux.SpotMin)

			fmt.Printf("  Windows:\n")
			fmt.Printf("    On-Demand: $%.4f\n", pricing.Windows.OnDemand)
			fmt.Printf("    Hybrid Benefit: $%.4f\n", pricing.Windows.HybridBenefit)
			fmt.Printf("    Basic: $%.4f\n", pricing.Windows.Basic)
			fmt.Printf("    Basic Spot: $%.4f\n", pricing.Windows.BasicSpot)
			fmt.Printf("    Spot Min: $%.4f\n", pricing.Windows.SpotMin)
		}

		// Print storage information if available
		if i.Storage != nil {
			fmt.Printf("\nStorage:\n")
			fmt.Printf("  Size: %d bytes\n", i.Storage.Size)
			fmt.Printf("  NVMe SSD: %t\n", i.Storage.NVMeSSD)
		}

		fmt.Println("-----------------------------")
	}
}

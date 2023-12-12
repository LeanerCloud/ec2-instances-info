package main

import (
	"fmt"
	"os"

	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info" // Adjust this import path as needed
)

func main() {
	data, err := ec2instancesinfo.OpenSearchData() // Call OpenSearchData instead of ElastiCacheData
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	instanceType := "c5.4xlarge.search" // specify the OpenSearch instance type you are interested in
	region := "us-east-1"               // specify the region

	for _, i := range *data {

		// fmt.Printf("Iterating over Type: %#v\n", i)
		if i.InstanceType == instanceType {
			if pricing, found := i.Pricing[region]; found {
				fmt.Println("Instance Type:", i.InstanceType)
				fmt.Println("Current Generation:", i.CurrentGeneration)
				fmt.Println("Instance Family:", i.InstanceFamily)
				fmt.Println("vCPU:", i.Vcpu)
				fmt.Println("Memory GiB:", i.MemoryGib) // Note: Now using MemoryGib
				fmt.Println("Service Code:", i.ServiceCode)
				fmt.Println("Service Name:", i.Servicename)

				// Print On-Demand pricing information
				fmt.Println("On-Demand cost:", pricing.OnDemand)

				// Print Reserved Plans pricing information
				fmt.Println("1-Year Standard Partial Upfront cost:", pricing.Reserved.YrTerm1StandardPartialUpfront)
				fmt.Println("3-Year Standard Partial Upfront cost:", pricing.Reserved.YrTerm3StandardPartialUpfront)
				fmt.Println("1-Year Standard No Upfront cost:", pricing.Reserved.YrTerm1StandardNoUpfront)
				fmt.Println("3-Year Standard No Upfront cost:", pricing.Reserved.YrTerm3StandardNoUpfront)
				// Add more fields as needed

				break // Stop iterating after finding the desired instance type
			}
		}
	}
}

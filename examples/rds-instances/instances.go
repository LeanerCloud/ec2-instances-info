package main

import (
	"fmt"
	"os"

	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info"
)

func main() {

	data, err := ec2instancesinfo.RDSData()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	instanceType := "db.t4g.medium"
	region := "us-east-1"

	for _, i := range *data {
		if i.InstanceType == instanceType {
			if pricing, found := i.Pricing[region]; found {
				fmt.Println("Instance type:", i.InstanceType)
				fmt.Println("CPU:", i.PhysicalProcessor)
				fmt.Println("Arch:", i.Arch)
				// Print other fields as needed

				fmt.Println("Aurora PostgreSQL OD cost:", pricing.AuroraPostgreSQL.OnDemand)
				fmt.Println("PostgreSQL OD cost:", pricing.PostgreSQL.OnDemand)
				// Print other pricing fields for different database engines

				break // Stop iterating after finding the desired instance type
			}
		}
	}

}

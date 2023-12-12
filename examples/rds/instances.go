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
				fmt.Println("Arch:", i.Arch)
				fmt.Println("Clock Speed:", i.ClockSpeed)
				fmt.Println("Current Generation:", i.CurrentGeneration)
				fmt.Println("Dedicated EBS Throughput:", i.DedicatedEbsThroughput)
				fmt.Println("EBS Baseline Bandwidth:", i.EbsBaselineBandwidth)
				fmt.Println("EBS Baseline IOPS:", i.EbsBaselineIops)
				fmt.Println("EBS Baseline Throughput:", i.EbsBaselineThroughput)
				fmt.Println("EBS IOPS:", i.EbsIops)
				fmt.Println("EBS Max Bandwidth:", i.EbsMaxBandwidth)
				fmt.Println("EBS Optimized:", i.EbsOptimized)
				fmt.Println("EBS Optimized By Default:", i.EbsOptimizedByDefault)
				fmt.Println("EBS Throughput:", i.EbsThroughput)
				fmt.Println("Enhanced Networking Supported:", i.EnhancedNetworkingSupported)
				fmt.Println("Family:", i.Family)
				fmt.Println("Instance Family:", i.InstanceFamily)
				fmt.Println("Instance Type Family:", i.InstanceTypeFamily)
				fmt.Println("Memory:", i.Memory)
				fmt.Println("Network Performance:", i.NetworkPerformance)
				fmt.Println("Physical Processor:", i.PhysicalProcessor)
				fmt.Println("Pretty Name:", i.PrettyName)

				// Print Pricing information for different database engines
				fmt.Println("Aurora PostgreSQL OD cost:", pricing.AuroraPostgreSQL.OnDemand)
				fmt.Println("PostgreSQL OD cost:", pricing.PostgreSQL.OnDemand)
				// ... and so on for other database engines

				fmt.Println("Processor Architecture:", i.ProcessorArchitecture)
				fmt.Println("Region Code:", i.RegionCode)
				// Print Regions information if needed
				// fmt.Println("Regions:", i.Regions)
				fmt.Println("Service Code:", i.Servicecode)
				fmt.Println("Service Name:", i.Servicename)
				fmt.Println("Storage:", i.Storage)
				fmt.Println("vCPU:", i.Vcpu)

				// Add more fields and pricing data as needed

				break // Stop iterating after finding the desired instance type
			}
		}
	}

}

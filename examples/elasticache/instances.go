package main

import (
	"fmt"
	"os"

	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info"
)

func main() {
	data, err := ec2instancesinfo.ElastiCacheData()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	instanceType := "cache.m6g.xlarge" // specify the instance type you are interested in
	region := "us-east-1"              // specify the region

	for _, i := range *data {

		// fmt.Printf("Iterating over Type: %#v\n", i)
		if i.InstanceType == instanceType {
			if pricing, found := i.Pricing[region]; found {
				fmt.Println("Instance Type:", i.InstanceType)
				fmt.Println("Current Generation:", i.CurrentGeneration)
				fmt.Println("Instance Family:", i.InstanceFamily)
				fmt.Println("vCPU:", i.Vcpu)
				fmt.Println("Memory:", i.Memory)
				fmt.Println("Network Performance:", i.NetworkPerformance)
				fmt.Println("Cache Engine:", i.CacheEngine)
				fmt.Println("Service Code:", i.ServiceCode)
				fmt.Println("Service Name:", i.Servicename)
				fmt.Println("Pretty Name:", i.PrettyName)

				// Print On-Demand pricing information
				fmt.Println("Memcached On-Demand cost:", pricing.Memcached.OnDemand)
				fmt.Println("Redis On-Demand cost:", pricing.Redis.OnDemand)

				// Print Reserved Plans pricing information
				fmt.Println("Memcached 1-Year Standard Partial Upfront cost:", pricing.Memcached.Reserved.YrTerm1StandardNoUpfront)
				fmt.Println("Memcached 3-Year Standard Partial Upfront cost:", pricing.Memcached.Reserved.YrTerm3StandardPartialUpfront)
				fmt.Println("Redis 1-Year Standard Partial Upfront cost:", pricing.Redis.Reserved.YrTerm1StandardPartialUpfront)
				fmt.Println("Redis 3-Year Standard Partial Upfront cost:", pricing.Redis.Reserved.YrTerm3StandardPartialUpfront)
				// Add more fields as needed

				break // Stop iterating after finding the desired instance type
			}
		}
	}
}

package ec2instancesinfo

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

type ElastiCacheInstanceData []ElastiCacheInstance

// PricingDetail represents the pricing information for each term.
type PricingDetail struct {
	OnDemand float64         `json:"ondemand"`
	Reserved ReservedPricing `json:"reserved"`
}

type ReservedPricing struct {
	YrTerm3StandardPartialUpfront float64 `json:"yrTerm3Standard.partialUpfront"`
	YrTerm1StandardPartialUpfront float64 `json:"yrTerm1Standard.partialUpfront"`
	YrTerm3StandardAllUpfront     float64 `json:"yrTerm3Standard.allUpfront"`
	YrTerm1StandardNoUpfront      float64 `json:"yrTerm1Standard.noUpfront"`
	YrTerm3StandardNoUpfront      float64 `json:"yrTerm3Standard.noUpfront"`
}

// ServicePricing represents the pricing for a particular service (e.g., Memcached, Redis).
type ServicePricing struct {
	Memcached PricingDetail `json:"Memcached"`
	Redis     PricingDetail `json:"Redis"`
}

// RegionPricing represents the pricing information for each region.
type RegionPricing struct {
	Region map[string]ServicePricing `json:"region"`
}

// ElastiCacheInstance represents the structure of each JSON object in the array.
type ElastiCacheInstance struct {
	ServiceCode                                    string                    `json:"servicecode"`
	InstanceType                                   string                    `json:"instanceType"`
	CurrentGeneration                              string                    `json:"currentGeneration"`
	InstanceFamily                                 string                    `json:"instanceFamily"`
	Vcpu                                           string                    `json:"vcpu"`
	Memory                                         string                    `json:"memory"`
	NetworkPerformance                             string                    `json:"networkPerformance"`
	CacheEngine                                    string                    `json:"cacheEngine"`
	RegionCode                                     string                    `json:"regionCode"`
	Servicename                                    string                    `json:"servicename"`
	Network_performance                            string                    `json:"network_performance"`
	Family                                         string                    `json:"family"`
	Instance_type                                  string                    `json:"instance_type"`
	Pricing                                        map[string]ServicePricing `json:"pricing"`
	Regions                                        map[string]string         `json:"regions"`
	PrettyName                                     string                    `json:"pretty_name"`
	Memcached1_6_MaxCacheMemory                    string                    `json:"memcached1.6-max_cache_memory"`
	Memcached1_6_NumThreads                        string                    `json:"memcached1.6-num_threads"`
	Redis6xClientOutputBufferLimitReplicaHardLimit string                    `json:"redis6.x-client-output-buffer-limit-replica-hard-limit"`
	Redis6xClientOutputBufferLimitReplicaSoftLimit string                    `json:"redis6.x-client-output-buffer-limit-replica-soft-limit"`
	Redis6xMaxmemory                               string                    `json:"redis6.x-maxmemory"`
	MaxClients                                     string                    `json:"max_clients"`
}

//go:embed data/elasticache-instances.json
var staticElastiCacheDataBody []byte

var elastiCacheDataBody, backupElastiCacheDataBody []byte

func ElastiCacheData() (*ElastiCacheInstanceData, error) {
	var d ElastiCacheInstanceData

	// Replace the handling of RDS data with ElastiCache data.
	if len(elastiCacheDataBody) > 0 {
		log.Println("We have updated ElastiCache data, trying to unmarshal it")
		err := json.Unmarshal(elastiCacheDataBody, &d)
		if err != nil {
			log.Printf("couldn't unmarshal the updated ElastiCache data, reverting to the backup ElastiCache data : %s", err.Error())
			err := json.Unmarshal(backupElastiCacheDataBody, &d)
			if err != nil {
				return nil, errors.Errorf("couldn't unmarshal backup ElastiCache data: %s", err.Error())
			}
			backupElastiCacheDataBody = []byte{}
		}
	} else {
		log.Println("Using the static ElastiCache instance type data")
		err := json.Unmarshal(staticElastiCacheDataBody, &d)
		if err != nil {
			return nil, errors.Errorf("couldn't unmarshal ElastiCache data: %s", err.Error())
		}
	}

	// Perform any ElastiCache-specific data processing here if needed.

	return &d, nil
}

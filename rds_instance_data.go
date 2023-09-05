package ec2instancesinfo

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

type RDSInstanceData []RDSInstance

type RDSInstance struct {
	Arch                        string `json:"arch,omitempty"`
	ClockSpeed                  string `json:"clockSpeed,omitempty"`
	CurrentGenerationRaw        string `json:"currentGeneration,omitempty"`
	CurrentGeneration           bool
	DedicatedEbsThroughput      string                     `json:"dedicatedEbsThroughput,omitempty"`
	EbsBaselineBandwidth        float64                    `json:"ebs_baseline_bandwidth,omitempty"`
	EbsBaselineIops             float64                    `json:"ebs_baseline_iops,omitempty"`
	EbsBaselineThroughput       float64                    `json:"ebs_baseline_throughput,omitempty"`
	EbsIops                     float64                    `json:"ebs_iops,omitempty"`
	EbsMaxBandwidth             float64                    `json:"ebs_max_bandwidth,omitempty"`
	EbsOptimized                bool                       `json:"ebs_optimized,omitempty"`
	EbsOptimizedByDefault       bool                       `json:"ebs_optimized_by_default,omitempty"`
	EbsThroughput               float64                    `json:"ebs_throughput,omitempty"`
	EnhancedNetworkingSupported string                     `json:"enhancedNetworkingSupported,omitempty"`
	Family                      string                     `json:"family,omitempty"`
	InstanceFamily              string                     `json:"instanceFamily,omitempty"`
	InstanceType                string                     `json:"instance_type,omitempty"`
	InstanceTypeFamily          string                     `json:"instanceTypeFamily,omitempty"`
	Memory                      float32                    `json:"memory,omitempty,string"`
	NetworkPerformance          string                     `json:"network_performance,omitempty"`
	NormalizationSizeFactor     string                     `json:"normalizationSizeFactor,omitempty"`
	PhysicalProcessor           string                     `json:"physicalProcessor,omitempty"`
	PrettyName                  string                     `json:"pretty_name,omitempty"`
	Pricing                     map[string]RDSRegionPrices `json:"pricing,omitempty"`
	ProcessorArchitecture       string                     `json:"processorArchitecture,omitempty"`
	RegionCode                  string                     `json:"regionCode,omitempty"`
	Regions                     Regions                    `json:"regions,omitempty"`
	Servicecode                 string                     `json:"servicecode,omitempty"`
	Servicename                 string                     `json:"servicename,omitempty"`
	Storage                     string                     `json:"storage,omitempty"`
	Vcpu                        int32                      `json:"vcpu,omitempty,string"`
}

type RDSRegionPrices struct {
	AuroraPostgreSQL             RDSPricing `json:"Aurora PostgreSQL,omitempty"`
	PostgreSQL                   RDSPricing `json:"PostgreSQL,omitempty"`
	SQLServer                    RDSPricing `json:"SQL Server,omitempty"`
	SQLServerOnPremiseForOutpost RDSPricing `json:"SQL Server (on-premise for Outpost),omitempty"`
	Oracle                       RDSPricing `json:"Oracle,omitempty"`
	MariaDB                      RDSPricing `json:"MariaDB,omitempty"`
	AuroraMySQL                  RDSPricing `json:"Aurora MySQL,omitempty"`
	MySQL                        RDSPricing `json:"MySQL,omitempty"`
}

type RDSPricing struct {
	OnDemand float64     `json:"ondemand"`
	Reserved RDSReserved `json:"reserved"`
}

type RDSReserved struct {
	StandardNoUpfront1Year          float64 `json:"yrTerm1Standard.noUpfront"`
	StandardNoUpfront3Years         float64 `json:"yrTerm3Standard.noUpfront"`
	StandardPartiallUpfront1Year    float64 `json:"yrTerm1Standard.partialUpfront"`
	StandardPartialUpfront3Years    float64 `json:"yrTerm3Standard.partialUpfront"`
	StandardAllUpfront1Year         float64 `json:"yrTerm1Standard.allUpfront"`
	StandardAllUpfront3Years        float64 `json:"yrTerm3Standard.allUpfront"`
	ConvertibleNoUpfront1Year       float64 `json:"yrTerm1Convertible.noUpfront"`
	ConvertibleNoUpfront3Years      float64 `json:"yrTerm3Convertible.noUpfront"`
	ConvertiblePartiallUpfront1Year float64 `json:"yrTerm1Convertible.partialUpfront"`
	ConvertiblePartialUpfront3Years float64 `json:"yrTerm3Convertible.partialUpfront"`
	ConvertibleAllUpfront1Year      float64 `json:"yrTerm1Convertible.allUpfront"`
	ConvertibleAllUpfront3Years     float64 `json:"yrTerm3Convertible.allUpfront"`
}

type Regions map[string]string

//go:embed data/rds-instances.json
var staticRDSDataBody []byte

var rdsDataBody, backupRDSDataBody []byte

func RDSData() (*RDSInstanceData, error) {
	var d RDSInstanceData

	// Similar to the EC2 instance data code, you can unmarshal your RDS data.
	// You will need to replace `staticRDSDataBody` with the appropriate RDS data.

	if len(rdsDataBody) > 0 {
		log.Println("We have updated RDS data, trying to unmarshal it")
		err := json.Unmarshal(rdsDataBody, &d)
		if err != nil {
			log.Printf("couldn't unmarshal the updated RDS data, reverting to the backup RDS data : %s", err.Error())
			err := json.Unmarshal(backupRDSDataBody, &d)
			if err != nil {
				return nil, errors.Errorf("couldn't unmarshal backup RDS data: %s", err.Error())
			}
			backupRDSDataBody = []byte{}
		}
	} else {
		log.Println("Using the static RDS instance type data")
		// Replace `staticRDSDataBody` with your RDS data variable name.
		err := json.Unmarshal(staticRDSDataBody, &d)
		if err != nil {
			return nil, errors.Errorf("couldn't unmarshal RDS data: %s", err.Error())
		}
	}

	// Similar to the EC2 instance data code, you can perform sorting and other operations here.

	for _, data := range d {
		if data.CurrentGenerationRaw == "Yes" {
			data.CurrentGeneration = true
		}
	}

	return &d, nil
}

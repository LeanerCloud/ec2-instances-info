package ec2instancesinfo

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

type OpenSearchInstanceData []OpenSearchInstance

type OpenSearchRegionPricing struct {
	OnDemand float64                   `json:"ondemand"`
	Reserved OpenSearchReservedPricing `json:"reserved"`
}

type OpenSearchReservedPricing struct {
	YrTerm3StandardPartialUpfront float64 `json:"yrTerm3Standard.partialUpfront"`
	YrTerm1StandardPartialUpfront float64 `json:"yrTerm1Standard.partialUpfront"`
	YrTerm3StandardAllUpfront     float64 `json:"yrTerm3Standard.allUpfront"`
	YrTerm1StandardNoUpfront      float64 `json:"yrTerm1Standard.noUpfront"`
	YrTerm3StandardNoUpfront      float64 `json:"yrTerm3Standard.noUpfront"`
}

type OpenSearchInstance struct {
	ServiceCode       string                             `json:"servicecode"`
	InstanceType      string                             `json:"instanceType"`
	CurrentGeneration string                             `json:"currentGeneration"`
	InstanceFamily    string                             `json:"instanceFamily"`
	Vcpu              int32                              `json:"vcpu,string"`
	MemoryGib         float64                            `json:"memoryGib,string"`
	RegionCode        string                             `json:"regionCode"`
	Servicename       string                             `json:"servicename"`
	Family            string                             `json:"family"`
	Instance_type     string                             `json:"instance_type"`
	Pricing           map[string]OpenSearchRegionPricing `json:"pricing"`
}

//go:embed data/opensearch-instances.json
var staticOpenSearchDataBody []byte

var openSearchDataBody, backupOpenSearchDataBody []byte

func OpenSearchData() (*OpenSearchInstanceData, error) {
	var d OpenSearchInstanceData

	// Replace the handling of RDS data with OpenSearch data.
	if len(openSearchDataBody) > 0 {
		log.Println("We have updated OpenSearch data, trying to unmarshal it")
		err := json.Unmarshal(openSearchDataBody, &d)
		if err != nil {
			log.Printf("couldn't unmarshal the updated OpenSearch data, reverting to the backup OpenSearch data : %s", err.Error())
			err := json.Unmarshal(backupOpenSearchDataBody, &d)
			if err != nil {
				return nil, errors.Errorf("couldn't unmarshal backup OpenSearch data: %s", err.Error())
			}
			backupOpenSearchDataBody = []byte{}
		}
	} else {
		log.Println("Using the static OpenSearch instance type data")
		err := json.Unmarshal(staticOpenSearchDataBody, &d)
		if err != nil {
			return nil, errors.Errorf("couldn't unmarshal OpenSearch data: %s", err.Error())
		}
	}

	// Perform any OpenSearch-specific data processing here if needed.

	return &d, nil
}

package ec2instancesinfo

// In this file we generate a raw data structure unmarshaled from the
// ec2instances.info JSON file, embedded into the binary at build time based on
// a generated snipped created using the go-bindata tool.

import (
	_ "embed"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// AWS Instances JSON Structure Definitions
type jsonInstance struct {
	Family             string          `json:"family"`
	EnhancedNetworking bool            `json:"enhanced_networking"`
	ECURaw             json.RawMessage `json:"ECU"`
	ECU                string
	VCPURaw            json.RawMessage `json:"vCPU"`
	VCPU               int
	PhysicalProcessor  string                  `json:"physical_processor"`
	Generation         string                  `json:"generation"`
	EBSIOPS            float32                 `json:"ebs_iops"`
	NetworkPerformance string                  `json:"network_performance"`
	EBSThroughput      float32                 `json:"ebs_throughput"`
	PrettyName         string                  `json:"pretty_name"`
	GPU                int                     `json:"GPU"`
	Pricing            map[string]RegionPrices `json:"pricing"`

	Storage *StorageConfiguration `json:"storage"`

	VPC struct {
		//    IPsPerENI int `json:"ips_per_eni"`
		//    MaxENIs   int `json:"max_enis"`
	} `json:"vpc"`

	Arch                     []string `json:"arch"`
	LinuxVirtualizationTypes []string `json:"linux_virtualization_types"`
	EBSOptimized             bool     `json:"ebs_optimized"`

	MaxBandwidth float32 `json:"max_bandwidth"`
	InstanceType string  `json:"instance_type"`

	Memory          float32
	MemoryRaw       json.RawMessage `json:"memory"`
	EBSMaxBandwidth float32         `json:"ebs_max_bandwidth"`
}

type StorageConfiguration struct {
	SSD     bool    `json:"ssd"`
	Devices int     `json:"devices"`
	Size    float32 `json:"size"`
	NVMeSSD bool    `json:"nvme_ssd"`
}

type RegionPrices struct {
	Linux              Pricing `json:"linux"`
	LinuxSQL           Pricing `json:"linuxSQL"`
	LinuxSQLEnterprise Pricing `json:"linuxSQLEnterprise"`
	LinuxSQLWeb        Pricing `json:"linuxSQLWeb"`
	MSWin              Pricing `json:"mswin"`
	MSWinSQL           Pricing `json:"mswinSQL"`
	MSWinSQLEnterprise Pricing `json:"mswinSQLEnterprise"`
	MSWinSQLWeb        Pricing `json:"mswinSQLWeb"`
	RHEL               Pricing `json:"rhel"`
	SLES               Pricing `json:"sles"`
	EBSSurcharge       float64 `json:"ebs,string"`
}

type Pricing struct {
	OnDemand float64  `json:"ondemand,string"`
	SpotMin  float64  `json:"spot_min,string"`
	SpotMax  float64  `json:"spot_max,string"`
	Reserved Reserved `json:"reserved"`
}

type Reserved struct {
	StandardNoUpfront1Year          float64 `json:"yrTerm1Standard.noUpfront,string"`
	StandardNoUpfront3Years         float64 `json:"yrTerm3Standard.noUpfront,string"`
	StandardPartiallUpfront1Year    float64 `json:"yrTerm1Standard.partialUpfront,string"`
	StandardPartialUpfront3Years    float64 `json:"yrTerm3Standard.partialUpfront,string"`
	StandardAllUpfront1Year         float64 `json:"yrTerm1Standard.allUpfront,string"`
	StandardAllUpfront3Years        float64 `json:"yrTerm3Standard.allUpfront,string"`
	ConvertibleNoUpfront1Year       float64 `json:"yrTerm1Convertible.noUpfront,string"`
	ConvertibleNoUpfront3Years      float64 `json:"yrTerm3Convertible.noUpfront,string"`
	ConvertiblePartiallUpfront1Year float64 `json:"yrTerm1Convertible.partialUpfront,string"`
	ConvertiblePartialUpfront3Years float64 `json:"yrTerm3Convertible.partialUpfront,string"`
	ConvertibleAllUpfront1Year      float64 `json:"yrTerm1Convertible.allUpfront,string"`
	ConvertibleAllUpfront3Years     float64 `json:"yrTerm3Convertible.allUpfront,string"`
}

//go:embed data/instances.json
var staticDataBody []byte

var dataBody, backupDataBody []byte

//------------------------------------------------------------------------------

// InstanceData is a large data structure containing pricing and specs
// information about all the EC2 instance types from all AWS regions.
type InstanceData []jsonInstance

// Data generates the InstanceData object based on data sourced from
// ec2instances.info. The data is available there as a JSON blob, which is
// converted into golang source-code by the go-bindata tool and unmarshaled into
// a golang data structure by this library.
func Data() (*InstanceData, error) {

	var d InstanceData

	if len(dataBody) > 0 {
		log.Println("We have updated data, trying to unmarshal it")
		err := json.Unmarshal(dataBody, &d)
		if err != nil {
			log.Printf("couldn't unmarshal the updated data, reverting to the backup data : %s", err.Error())
			err := json.Unmarshal(backupDataBody, &d)
			if err != nil {
				return nil, errors.Errorf("couldn't unmarshal backup data: %s", err.Error())
			}
			backupDataBody = []byte{}
		}
	} else {
		log.Println("Using the static instance type data")
		err := json.Unmarshal(staticDataBody, &d)
		if err != nil {
			return nil, errors.Errorf("couldn't unmarshal data: %s", err.Error())
		}

	}

	// Handle "N/A" values in the VCPU field for i3.metal instance type and
	// string ("variable") and integer values in the ECU field
	for i := range d {
		var vcpu, intECU int
		var stringECU string
		var memory float32

		if err := json.Unmarshal(d[i].VCPURaw, &vcpu); err == nil {
			d[i].VCPU = vcpu
		}

		if err := json.Unmarshal(d[i].MemoryRaw, &memory); err == nil {
			d[i].Memory = memory
		}

		if err := json.Unmarshal(d[i].ECURaw, &intECU); err == nil {
			d[i].ECU = strconv.Itoa(intECU)
		} else if err = json.Unmarshal(d[i].ECURaw, &stringECU); err == nil {
			d[i].ECU = stringECU
		}
	}

	sort.Slice(d, func(i, j int) bool {
		// extract the instance family, such as "c5" for "c5.large"
		family_i := strings.Split(d[i].InstanceType, ".")[0]
		family_j := strings.Split(d[j].InstanceType, ".")[0]

		// we first compare only based on the family
		switch strings.Compare(family_i, family_j) {
		case -1:
			return true
		case 1:
			return false
		}

		// within the same family we compare by memory size, but always keeping metal instances last
		return d[i].Memory < d[j].Memory || strings.HasSuffix(d[j].InstanceType, "metal")
	})

	return &d, nil
}

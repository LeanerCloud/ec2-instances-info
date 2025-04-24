// azure_vm.go with updated fields
package ec2instancesinfo

// AzureInstanceData represents data about a specific Azure VM instance type
type AzureInstanceData struct {
	ACU                   int                     `json:"ACU"`
	GPU                   string                  `json:"GPU"`
	AcceleratedNetworking bool                    `json:"accelerated_networking"`
	Arch                  []string                `json:"arch"`
	AvailabilityZones     map[string]interface{}  `json:"availability_zones"`
	CachedDisk            int                     `json:"cached_disk"`
	CapacitySupport       bool                    `json:"capacity_support"`
	Category              string                  `json:"category"`
	Confidential          interface{}             `json:"confidential"` // Can be false or "SNP"
	Encryption            bool                    `json:"encryption"`
	Family                string                  `json:"family"`
	Hibernation           interface{}             `json:"hibernation"`
	HypervGenerations     string                  `json:"hyperv_generations"`
	InstanceType          string                  `json:"instance_type"`
	IOPS                  interface{}             `json:"iops"`
	LowPriority           bool                    `json:"low_priority"`
	Memory                float64                 `json:"memory"`
	MemoryMaintenance     bool                    `json:"memory_maintenance"`
	PremiumIO             bool                    `json:"premium_io"`
	PrettyName            string                  `json:"pretty_name"`
	PrettyNameAzure       string                  `json:"pretty_name_azure"`
	Pricing               map[string]AzurePricing `json:"pricing"`
	RDMA                  bool                    `json:"rdma"`
	ReadIO                int                     `json:"read_io"`
	Size                  int                     `json:"size"`
	Storage               *AzureStorage           `json:"storage"`
	TrustedLaunch         interface{}             `json:"trusted_launch"`
	UltraSSD              bool                    `json:"ultra_ssd"`
	UncachedDisk          int                     `json:"uncached_disk"`
	UncachedDiskIO        int                     `json:"uncached_disk_io"`
	VCPU                  int                     `json:"vcpu"`
	VCPUsAvailable        int                     `json:"vcpus_available"`
	VCPUsPerCore          int                     `json:"vcpus_percore"`
	VMDeployment          string                  `json:"vm_deployment"`
	WriteIO               int                     `json:"write_io"`
}

// AzureStorage contains Azure VM storage information
type AzureStorage struct {
	Devices       interface{} `json:"devices"`
	MaxWriteDisks interface{} `json:"max_write_disks"`
	NVMeSSD       interface{} `json:"nvme_ssd"` // Can be string with numeric value or false
	Size          int64       `json:"size"`
}

// AzurePricing contains pricing information for a region
type AzurePricing struct {
	Linux   AzureOSPricing `json:"linux"`
	Windows AzureOSPricing `json:"windows"`
}

// AzureOSPricing contains pricing details for a specific OS
type AzureOSPricing struct {
	Basic       float64                `json:"basic"`
	BasicSpot   float64                `json:"basic-spot"`
	LowPriority float64                `json:"lowpriority,omitempty"`
	OnDemand    float64                `json:"ondemand"`
	Reserved    map[string]interface{} `json:"reserved"`
	SpotMin     float64                `json:"spot_min"`
	// Windows specific
	HybridBenefit float64 `json:"hybridbenefit,omitempty"`
}

// Global variables for Azure VM data
var (
	azureDataBody       []byte
	azureBackupDataBody []byte
	azureStaticDataBody = []byte(`[]`) // Will be replaced with actual embedded data
	azureInstancesData  *[]AzureInstanceData
)

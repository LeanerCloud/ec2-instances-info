package ec2instancesinfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultAzureAPIHost = "api.leanercloud.com"
	defaultAzureAPIPath = "/azure-vm-info/"
)

// AzureData returns the Azure VM instance data, either from the in-memory cache
// or by loading the data from the JSON file
func AzureData() (*[]AzureInstanceData, error) {
	if azureInstancesData != nil {
		return azureInstancesData, nil
	}

	var jsonData []byte
	if len(azureDataBody) > 0 {
		jsonData = azureDataBody
	} else {
		// Try to load from the data directory
		data, err := ioutil.ReadFile("data/azure-instances.json")
		if err != nil {
			// If file not found, use static embedded data
			jsonData = azureStaticDataBody
		} else {
			jsonData = data
		}
	}

	instances := []AzureInstanceData{}
	err := json.Unmarshal(jsonData, &instances)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal Azure VM data")
	}

	azureInstancesData = &instances
	return azureInstancesData, nil
}

// UpdateAzureData updates the Azure VM instance data by downloading it from the API
func UpdateAzureData(apiHost, apiKey *string) error {
	if apiHost == nil {
		s := defaultAzureAPIHost
		apiHost = &s
	}

	log.Printf("Dynamic Azure data size before: %d, downloading new Azure VM data.", len(azureDataBody))

	url, err := getAzureDataURL(*apiHost, *apiKey)
	if err != nil {
		return fmt.Errorf("error getting Azure data download URL: %s", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return fmt.Errorf("error creating Azure data request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Azure data request: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading Azure data response body: %s", err)
	}

	if len(azureDataBody) > 0 {
		azureBackupDataBody = azureDataBody
	} else {
		azureBackupDataBody = azureStaticDataBody
	}

	azureDataBody = body
	log.Println("Azure data size after:", len(azureDataBody))

	// Reset the in-memory instance data so it will be reloaded next time it's accessed
	azureInstancesData = nil

	// Optionally write a local copy for debugging
	out, err := os.Create("azure_instances_dump.json")
	if err != nil {
		log.Println("Couldn't create Azure dump file:", err)
	} else {
		defer out.Close()
		out.Write(body)
	}

	return nil
}

func getAzureDataURL(apiHost, apiKey string) (*string, error) {
	url := fmt.Sprintf("https://%s%s", apiHost, defaultAzureAPIPath)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating Azure request: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending Azure request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected HTTP status for Azure data: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Azure response body: %s", err)
	}
	ret := string(body)

	return &ret, nil
}

// AzureUpdater continuously updates the Azure VM instance data at the specified interval
func AzureUpdater(refreshDays int, apiHost, apiKey *string) error {
	if apiKey == nil {
		log.Println("Azure API key is missing")
		return fmt.Errorf("Azure API key is missing")
	}

	if apiHost == nil {
		host := defaultAzureAPIHost
		apiHost = &host
	}

	if refreshDays <= 0 {
		refreshDays = defaultRefreshInterval
	}
	refreshInterval := time.Duration(refreshDays) * 24 * time.Hour

	if err := UpdateAzureData(apiHost, apiKey); err != nil {
		log.Printf("Failed to download updated Azure VM data: %s", err.Error())
		return fmt.Errorf("error downloading new Azure VM data: %s", err)
	}

	// refresh the file every refreshInterval
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := UpdateAzureData(apiHost, apiKey)
			if err != nil {
				log.Println("Error refreshing Azure VM data:", err)
			}
		}
	}
}
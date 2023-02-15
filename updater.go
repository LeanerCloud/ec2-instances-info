package ec2instancesinfo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultAPIHost         = "api.leanercloud.com"
	defaultAPIPath         = "/instance-type-info/"
	defaultRefreshInterval = 7
)

func getDataURL(apiHost, apiKey string) (*string, error) {
	url := fmt.Sprintf("https://%s%s", apiHost, defaultAPIPath)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected HTTP status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}
	ret := string(body)

	return &ret, nil
}

func UpdateData(apiHost, apiKey *string) error {
	if apiHost == nil {
		s := defaultAPIHost
		apiHost = &s
	}

	log.Printf("Dynamic data size before: %d, downloading new instance type data.", len(dataBody))

	url, err := getDataURL(*apiHost, *apiKey)

	if err != nil {
		return fmt.Errorf("error getting download URL file: %s", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}
	if len(dataBody) > 0 {
		backupDataBody = dataBody
	} else {
		backupDataBody = staticDataBody
	}

	dataBody = body
	log.Println("Data size after:", len(dataBody))

	out, err := os.Create("instances_dump.json")
	if err != nil {
		panic("Couldn't create dump file")
	}

	defer out.Close()
	out.Write(body)

	return nil
}

func Updater(refreshDays int, apiHost, apiKey *string) error {
	if apiKey == nil {
		log.Println("API key is missing")
		return fmt.Errorf("API key is missing")
	}

	if apiHost == nil {
		host := defaultAPIHost
		apiHost = &host
	}

	if refreshDays <= 0 {
		refreshDays = defaultRefreshInterval
	}
	refreshInterval := time.Duration(refreshDays) * 24 * time.Hour

	if err := UpdateData(apiHost, apiKey); err != nil {
		log.Printf("Failed to download updated data: %s", err.Error())
		return fmt.Errorf("error downloading new data: %s", err)
	}

	// refresh the file every refreshInterval
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := UpdateData(apiHost, apiKey)
			if err != nil {
				log.Println("Error refreshing data:", err)
			}
		}
	}
}

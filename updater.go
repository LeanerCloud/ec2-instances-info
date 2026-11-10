package ec2instancesinfo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	defaultAPIHost         = "api.leanercloud.com"
	defaultAPIPath         = "/instance-type-info/"
	defaultRefreshInterval = 7
)

// httpClient is replaced in tests.
var httpClient = &http.Client{}

// withoutURL strips the request URL that net/http embeds in transport errors.
// The data file URL returned by the API is a signed link, so it must not end
// up in error text that callers log.
func withoutURL(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) {
		return ue.Err
	}
	return err
}

func getDataURL(apiHost, apiKey string) (*string, error) {
	url := fmt.Sprintf("https://%s%s", apiHost, defaultAPIPath)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := httpClient.Do(req)
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

	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", withoutURL(err))
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %s", withoutURL(err))
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received unexpected HTTP status downloading the data file: %s", resp.Status)
	}

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

	// The dump is a debugging aid and the working directory is read-only on
	// some platforms (e.g. Lambda), so failing to write it must not fail an
	// update that already succeeded.
	out, err := os.Create("instances_dump.json")
	if err != nil {
		log.Printf("Couldn't create the dump file, skipping it: %s", err)
		return nil
	}

	defer out.Close()
	if _, err := out.Write(body); err != nil {
		log.Printf("Couldn't write the dump file: %s", err)
	}

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

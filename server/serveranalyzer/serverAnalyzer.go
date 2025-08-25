package serveranalyzer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ip checks will be done on the server, you get packet metadata here so just ingest and call the apis async with golang.

var baseURL = "https://api.abuseipdb.com/api/v2/check"

type abuseIPDBResponse struct {
	Data struct {
		AbuseConfidenceScore int `json:"abuseConfidenceScore"`
	} `json:"data"`
}

func IpCheckAbuseIPDB(ip string, apiKey *string, ipdbClient *http.Client) (int, error) {
	params := url.Values{}
	params.Add("ipAddress", ip)

	fullURL := baseURL + "?" + params.Encode()

	score := 0

	req, err := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Key", *apiKey)
	req.Header.Add("Accept", "application/json")

	if err != nil {
		return score, err
	}

	resp, err := ipdbClient.Do(req)

	if err != nil {
		return score, err
	}

	if resp.StatusCode != 200 {
		return score, fmt.Errorf("received non-200 status code for ipdb lookup")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return score, err
	}

	var parsed abuseIPDBResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return score, err
	}

	score = parsed.Data.AbuseConfidenceScore
	return score, nil

}

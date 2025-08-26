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

	fmt.Println("before first err")
	if err != nil {
		return score, err
	}

	resp, err := ipdbClient.Do(req)

	fmt.Println("before second err")

	if err != nil {
		fmt.Println("why are we returning in here? ", score, err)
		return score, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return score, fmt.Errorf("received non-200 status code for ipdb lookup")
	}

	fmt.Println("was 200")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return score, err
	}

	var parsed abuseIPDBResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return score, err
	}
	fmt.Println("could marshal")

	score = parsed.Data.AbuseConfidenceScore
	fmt.Println("score was: ", score)
	return score, nil

}

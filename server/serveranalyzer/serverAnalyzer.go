package serveranalyzer

import (
	"net/http"
	"net/url"
)

// ip checks will be done on the server, you get packet metadata here so just ingest and call the apis async with golang.

var baseURL = "https://api.abuseipdb.com/api/v2/check"

func ipCheckAbuseIPDB(ip string) {
	params := url.Values{}
	params.Add("ipAddress", ip)

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL)

}

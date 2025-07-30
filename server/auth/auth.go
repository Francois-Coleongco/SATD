package auth

import (
	"SATD/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func AuthToDash(dashClient *http.Client, attempts int, dashboardServerAuthAddr string, dashUserCreds types.DashCreds, dashboardJWT *string) error {

	log.Println("this is user in authtodash ", dashUserCreds.Username)
	log.Println("this is user in authtodash ", dashUserCreds.Username)
	log.Println("this is pass in authtodash ", dashUserCreds.Password)
	log.Println("this is pass in authtodash ", dashUserCreds.Password)

	creds, err := json.Marshal(dashUserCreds)

	if err != nil {
		log.Printf("couldn't marshal dashUserCreds, error thrown: %s\n", err)
	}

	reader := bytes.NewReader(creds)

	success := false

	for range attempts {

		req, err := http.NewRequest("POST", dashboardServerAuthAddr, reader)
		req.Header.Set("Content-Type", "application/json")
		reader.Seek(0, io.SeekStart)

		if err != nil {
			log.Printf("failed authToDash httpRequest, error thrown: %s\n", err)
			continue
		}

		res, err := dashClient.Do(req)

		if err != nil {
			log.Printf("couldn't send request in authToDash, error thrown: %s\n", err)
			continue
		}

		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)

		if err != nil {
			log.Printf("unable to readAll from res.Body in authToDash, error thrown: %s\n", err)
		}

		if res.StatusCode == 200 {
			success = true
			fmt.Println("DATA FROM AUTH ENDPOINT WAS: ")

			var jwt types.JWT
			err := json.Unmarshal(data, &jwt)

			if err != nil {
				log.Printf("error unmarshalling jwt token, error thrown: %s", err)
			}

			*dashboardJWT = jwt.Token

			break
		}

	}

	if !success {
		return fmt.Errorf("couldn't authenticate to node server after all attempts")
	}

	return nil
}

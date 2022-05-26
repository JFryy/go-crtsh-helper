package crtsh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Cert struct matching json response of crt.sh api.
type Cert struct {
	IssuerCaID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
}

// Get certificate response from crt.sh and returns marshalled json as list of structs.
func getResponse(url string) ([]Cert, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not make request to crt.sh - %w", err)
	}
	defer resp.Body.Close() // no error checking, as if this refuses the program is just dead
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body of response - %w", err)
	}
	var certs []Cert
	err = json.Unmarshal(body, &certs)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response - %w", err)
	}
	return certs, nil
}

// GetCerts by domain name, organization name or certificate fingerprint.
// Input parameter as 'string' for single query, input as slice for joining multiple queries/responses.
func GetCerts(query interface{}) ([]Cert, error) {
	baseURL := "https://crt.sh/"
	output := "json"
	var allcerts [][]Cert

	switch query.(type) {
	case []string:
		s := reflect.ValueOf(query)
		for i := 0; i < s.Len(); i++ {
			url := fmt.Sprintf("%s?q=%s&output=%s", baseURL, s.Index(i), output)
			certs, err := getResponse(url)
			if err != nil {
				return nil, fmt.Errorf("could not get certs - %w", err)
			}
			allcerts = append(allcerts, certs)
		}

	case string:
		url := fmt.Sprintf("%s?q=%s&output=%s", baseURL, query, output)
		certs, err := getResponse(url)
		if err != nil {
			return nil, fmt.Errorf("could not get certs - %w", err)
		}
		return certs, nil


	default:
		return nil, fmt.Errorf("unaccepted argument for query")
	}


	var mergedCerts []Cert
	for host := range allcerts {
		for certs := range allcerts[host] {
			mergedCerts = append(mergedCerts, allcerts[host][certs])
		}
	}
	return mergedCerts, nil
}

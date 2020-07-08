package crtsh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
}

// Get certificate response from crt.sh and returns marshalled json as list of structs.
func getResponse(url string) []Cert {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var certs []Cert
	json.Unmarshal(body, &certs)
	return certs
}

// GetCerts by domain name, organization name or certificate fingerprint.
// Input parameter as 'string' for single query, input as slice for joining multiple queries/responses.
func GetCerts(query interface{}) []Cert {
	baseURL := "https://crt.sh/"
	output := "json"
	var certs []Cert
	var allcerts [][]Cert

	switch query.(type) {
	case []string:
		s := reflect.ValueOf(query)
		for i := 0; i < s.Len(); i++ {
			url := fmt.Sprintf("%s?q=%s&output=%s", baseURL, s.Index(i), output)
			certs = getResponse(url)
			allcerts = append(allcerts, certs)
		}

	case string:
		url := fmt.Sprintf("%s?q=%s&output=%s", baseURL, query, output)
		certs = getResponse(url)

	default:
		log.Fatal("Unaccepted argument for query.")
	}

	if len(allcerts) == 0 {
		return certs
	}

	var mergedCerts []Cert
	for host := range allcerts {
		for certs := range allcerts[host] {
			mergedCerts = append(mergedCerts, allcerts[host][certs])
		}
	}
	return mergedCerts
}

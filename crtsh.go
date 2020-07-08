package crtsh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

// Struct matching json response of getcert api.
type Cert struct {
	IssuerCaID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
}

// Get certificate response from crt.sh and returns marshal struct
func GetCerts(url string) []Cert {
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

// Query by domain name, organization name or certificate fingerprint.
// Input parameter has hand
func query(query interface{}) []Cert {
	base_url := "https://crt.sh/"
	output := "json"
	var certs []Cert
	var allcerts [][]Cert

	switch query.(type) {
	case []string:
		s := reflect.ValueOf(query)
		for i := 0; i < s.Len(); i++ {
			url := fmt.Sprintf("%s?q=%s&output=%s", base_url, s.Index(i), output)
			certs = GetCerts(url)
			allcerts = append(allcerts, certs)
		}

	case string:
		url := fmt.Sprintf("%s?q=%s&output=%s", base_url, query, output)
		certs = GetCerts(url)

	default:
		log.Fatal("Unaccepted argument for query.")
	}

	if len(allcerts) == 0 {
		return certs
	}

	var merged_certs []Cert
	for host := range allcerts {
		for certs := range allcerts[host] {
			merged_certs = append(merged_certs, allcerts[host][certs])
		}

	}
	return merged_certs
}

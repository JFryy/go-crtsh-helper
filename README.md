## crt.sh Golang Client
This is a _tiny_ and not feature-full http client in golang for retrieving certificate data from the website `crt.sh`, 
which is a web interface (Run by Comodo) for the distributed database known as the `certificate transparency logs`.
Please see ['What is Certificate Transparency?'](certificate-transparency.org/what-is-ct) for more information.

### Installation
`go get -u github.com/JFryy/go-crtsh-api`

### Usage

Import the package
```go
import (
	"fmt"
	"github.com/JFryy/go-crtsh-api"
)
```

Use the GetCerts method for a single organization/domain
```go
func main() {
	// Print Issuer and NameValue certificates for a domain and subdomains
	certs := crtsh.GetCerts("google.com")
	for i := range certs {
		fmt.Println(certs[i].IssuerName, certs[i].NameValue)
	}
}
```

Or use the GetCerts method for multiple organizations/domains
```go
func main() {
    domains := []string{"google.com", "reddit.com"}
	certs := crtsh.GetCerts(domains)
	for i := range certs {
		fmt.Println(certs[i].IssuerName, certs[i].NameValue)
	}
}
```

The cert struct and json from `crt.sh` has the following attributes:
```text
IssuerCaID     int
IssuerName     string
NameValue      string
ID             int64
EntryTimestamp string
NotBefore      string
NotAfter       string
```


#### Notes
_Upstream changes can and will break functionality of this package. Please use at your own risk._

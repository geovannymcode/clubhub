package command

import (
	"log"
	"net"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func RunWhoIs(ipAddr string) map[string][]string {
	ipObj := net.ParseIP(ipAddr)
	if ipObj == nil {
		log.Println("Invalid IP Address!")
		return make(map[string][]string)
	}

	result, err := whois.Whois(ipAddr)
	if err != nil {
		log.Println("Error fetching WHOIS data:", err)
		return nil
	}

	parsedResult, err := whoisparser.Parse(result)
	if err != nil {
		log.Println("Error parsing WHOIS data:", err)
		return nil
	}

	var outPut map[string][]string
	outPut = make(map[string][]string)

	if parsedResult.Domain != nil {
		outPut["DomainName"] = []string{parsedResult.Domain.Name}
		outPut["CreationDate"] = []string{parsedResult.Domain.CreatedDate}
		outPut["ExpirationDate"] = []string{parsedResult.Domain.ExpirationDate}
	}
	if parsedResult.Registrar != nil {
		outPut["RegistrarName"] = []string{parsedResult.Registrar.Name}
		outPut["RegistrarEmail"] = []string{parsedResult.Registrar.Email}
	}

	return outPut
}

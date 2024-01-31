package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
)

func GetDataSSl(address string) (coreDomain.SSL, error) {
	response, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + address)

	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Printf("status code error: %d %s", response.StatusCode, response.Status)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var responseSSL coreDomain.SSL
	json.Unmarshal(responseData, &responseSSL)

	return responseSSL, err
}

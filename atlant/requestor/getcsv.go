package requestor

import (
	aerr "atlant/errors"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetCsvFile(url string, timeoutSec uint64) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeoutSec),
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error during downloading file: %v", err)
		return nil, aerr.NewServiceError(aerr.ResourceUnavailable)
	}
	if resp.StatusCode != 200 {
		log.Printf("External server returns non 200: status %s", http.StatusText(resp.StatusCode))
		return nil, aerr.NewServiceError(aerr.ResourceUnavailable)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error during reading body response")
		return nil, errors.New("Internal error: can't read response")
	}
	return result, nil
}

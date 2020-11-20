package requestor

import (
	"net/http"
	"time"
	"errors"
	"log"
	"io/ioutil"
)

func GetCsvFile(url string, timeoutSec uint64) ([]byte, error) {
	client := http.Client {
		Timeout : time.Second * time.Duration(timeoutSec),
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error during downloading file: %v", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("External server returns non 200: status %s", http.StatusText(resp.StatusCode))
		return nil, errors.New("Can't download csv file")
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error during reading body response")
		return nil, errors.New("Internal error: can't read response")
	}
	return result, nil
}

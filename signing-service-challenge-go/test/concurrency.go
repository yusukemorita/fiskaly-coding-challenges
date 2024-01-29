package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/google/uuid"
)

func main() {
	deviceID := uuid.NewString()

	log.Printf("creating device (id: %s)", deviceID)
	_, err := sendJsonRequest(
		http.MethodPost,
		"http://localhost:8080/api/v0/signature_devices",
		api.CreateSignatureDeviceRequest{
			ID:        deviceID,
			Label:     "my rsa key",
			Algorithm: "RSA",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	send100Requests := func(name string) {
		for i := 0; i < 100; i++ {
			log.Printf("goroutine %s: sending request %d", name, i)
			_, err = sendJsonRequest(
				http.MethodPost,
				fmt.Sprintf("http://localhost:8080/api/v0/signature_devices/%s/signatures", deviceID),
				api.SignTransactionRequest{
					DataToBeSigned: fmt.Sprintf("some-data-%s", i),
				},
			)
			if err != nil {
				log.Fatal(err)
			}
		}
		wg.Done()
	}

	wg.Add(2)
	go send100Requests("A")
	go send100Requests("B")
	wg.Wait()
}

func sendJsonRequest(
	httpMethod string,
	url string,
	serializableData ...any,
) (*http.Response, error) {
	var bodyReader io.Reader
	if len(serializableData) > 0 {
		jsonBytes, err := json.Marshal(serializableData[0])
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	request, err := http.NewRequest(httpMethod, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %s", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func readBody(response *http.Response) string {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	return string(body)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Vehicle structure must match JSON schema
type Vehicle struct {
	Id    int                    `json:"id"`
	Year  int                    `json:"year"`
	Brand map[string]interface{} `json:"brand"`
}

// Vehicles structure must match JSON schema
type Vehicles struct {
	Vehicle []Vehicle `json:"vehicles"`
}

// Data structure must match JSON schema
type Data struct {
	Vehicles Vehicles `json:"data"`
}

func output(w http.ResponseWriter, r *http.Request) {
	// setup http client
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.maxposter.ru/partners-api/vehicles/active", nil)
	// add params
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic Y2hlc3RhdnRvQG1heHBvc3Rlci5ydTpuNWsxZzdxRA")
	if err != nil {
		log.Fatal(err)
	}
	// perform request
	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Fatal(respErr)
	}

	// create a []byte to Unmarshal it into the struct
	answer, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr, " read err")
	}

	// write []byte to the struct
	var data Data
	unmarhalErr := json.Unmarshal(answer, &data)
	if unmarhalErr != nil {
		log.Fatal(unmarhalErr, " unmarshal err")
	}

	_, writeErr := fmt.Fprint(w, data)
	if writeErr != nil {
		log.Fatal(writeErr, " writing to page err")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
}

func main() {
	http.HandleFunc("/", output)
	serverErr := http.ListenAndServe(":8080", nil)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}

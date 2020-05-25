package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch"
)

var es *elasticsearch.Client

func init() {
	es = createESClient()
}

func main() {
	log.SetFlags(0)

	var (
		r map[string]interface{}
		// wg sync.WaitGroup
	)

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	// log.Printf("Client: %s", )
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

}

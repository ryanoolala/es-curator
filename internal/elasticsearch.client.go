package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/elastic/go-elasticsearch"
	awsgosignv4 "github.com/jriquelme/awsgosigv4"
)

// Cluster configures the cluster
//
type Cluster struct {
	Client *elasticsearch.Client
	ESType string
}

// type Snapshot struct {
// 	client    *elasticsearch.Client
// 	IndexName string
// }

// func newSnapshot(c ClusterConfig) (*Snapshot, error) {

// }

// Creates elasticsearch client and returns cluster object
func createESClient(configuration Configurations) *Cluster {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	esCfg := elasticsearch.Config{
		Addresses: []string{
			configuration.ESHost,
		},

		Transport: &awsgosignv4.SignV4SDKV2{
			RoundTripper: http.DefaultTransport,
			Credentials:  cfg.Credentials,
			Region:       cfg.Region,
			Service:      "es",
			Now:          time.Now,
		},
	}
	client, err := elasticsearch.NewClient(esCfg)

	cluster := Cluster{Client: client, ESType: "aws"}

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	return &cluster
}

// Call cluster endpoint "/" to retrieve information on the cluster
//
func info(c *Cluster) (*Info, error) {
	var info Info
	res, err := c.Client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Printf("Server: %s", info.Version.Number)
	log.Println(strings.Repeat("~", 37))

	return &info, err
}

// Info "/" response json
//
type Info struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
	Version     struct {
		Number                           string    `json:"number"`
		BuildFlavor                      string    `json:"build_flavor"`
		BuildType                        string    `json:"build_type"`
		BuildHash                        string    `json:"build_hash"`
		BuildDate                        time.Time `json:"build_date"`
		BuildSnapshot                    bool      `json:"build_snapshot"`
		LuceneVersion                    string    `json:"lucene_version"`
		MinimumWireCompatibilityVersion  string    `json:"minimum_wire_compatibility_version"`
		MinimumIndexCompatibilityVersion string    `json:"minimum_index_compatibility_version"`
	} `json:"version"`
	Tagline string `json:"tagline"`
}

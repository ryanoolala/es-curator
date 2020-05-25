package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch"
)

var (
	fixtures = make(map[string]io.ReadCloser)
)

// Initialise test fixtures
//
func init() {
	fixtureFiles, err := filepath.Glob("testdata/*.json")
	if err != nil {
		panic(fmt.Sprintf("Cannot glob fixture files: %s", err))
	}

	for _, fpath := range fixtureFiles {
		f, err := ioutil.ReadFile(fpath)
		if err != nil {
			panic(fmt.Sprintf("Cannot read fixture file: %s", err))
		}

		fixtures[filepath.Base(fpath)] = ioutil.NopCloser(bytes.NewReader(f))
	}
}

// fixture function to read json
//
func fixture(fname string) io.ReadCloser {
	out := new(bytes.Buffer)
	b1 := bytes.NewBuffer([]byte{})
	b2 := bytes.NewBuffer([]byte{})
	tr := io.TeeReader(fixtures[fname], b1)

	defer func() { fixtures[fname] = ioutil.NopCloser(b1) }()
	io.Copy(b2, tr)
	out.ReadFrom(b2)

	return ioutil.NopCloser(out)
}

// MockTransport for elasticsearch response
//
type MockTransport struct {
	Response    *http.Response
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}

// Test Info function
//
func TestInfo(t *testing.T) {
	t.Parallel()

	mocktrans := MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{}`)),
		},
	}

	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Transport: &mocktrans,
	})

	cluster := &Cluster{Client: client}

	if err != nil {
		t.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	t.Run("Info", func(t *testing.T) {
		mocktrans.Response = &http.Response{
			StatusCode: http.StatusOK,
			Body:       fixture("info.json"),
		}

		res, err := info(cluster)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if res.Version.Number != "6.7.0" {
			t.Errorf("Unexpected version number, want=6.7.0, got=%s", res.Version.Number)
		}
	})
}

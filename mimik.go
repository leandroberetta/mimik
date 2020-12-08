package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type endpoint struct {
	Name        string       `json:"name"`
	Path        string       `json:"path"`
	Method      string       `json:"method"`
	Connections []connection `json:"connections"`
}

type connection struct {
	Name   string `json:"name"`
	Port   string `json:"port"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type service struct {
	Name      string
	Port      string
	Endpoints []endpoint
}

type response struct {
	Name             string     `json:"name"`
	Version          string     `json:"version"`
	Path             string     `json:"path"`
	StatusCode       int        `json:"statusCode"`
	UpstreamResponse []response `json:"upstreamResponse"`
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func main() {
	service, _ := newService(
		os.Getenv("MIMIK_SERVICE_NAME"),
		os.Getenv("MIMIK_SERVICE_PORT"),
		os.Getenv("MIMIK_ENDPOINTS_FILE"),
		getVersion(os.Getenv("MIMIK_LABELS_FILE")))
	client := &http.Client{}
	http.HandleFunc("/", endpointHandler(service, client))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", service.Port), nil))
}

func getVersion(fileName string) string {
	version := "v1"
	labelsFile, err := os.Open(fileName)
	if err != nil {
		return version
	}
	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), "=")
		if values[0] == "version" {
			version = values[1]
		}
	}
	return version
}

func newService(name, port, fileName, version string) (service, error) {
	service := service{Name: name, Port: port}
	err := loadEndpoints(fileName, &service.Endpoints)
	return service, err
}

func loadEndpoints(fileName string, endpoints *[]endpoint) error {
	file, err := os.Open(fileName)
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, endpoints)
	return err
}

func endpointHandler(service service, client httpClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := response{Name: service.Name, Version: "v1", StatusCode: http.StatusNotFound}
		ch := make(chan response)
		headers := getHeaders(r.Header)
		for _, endpoint := range service.Endpoints {
			resp.Path = r.URL.Path
			if endpoint.Path == r.URL.Path {
				resp.StatusCode = http.StatusOK
				upstreamResponse := make([]response, len(endpoint.Connections))
				for _, connection := range endpoint.Connections {
					go handleReq(makeURL(connection), connection.Method, headers, client, ch)
				}
				for i := range endpoint.Connections {
					upstreamResponse[i] = <-ch
				}
				resp.UpstreamResponse = upstreamResponse
			}
		}
		responseJSON, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

func getHeaders(header http.Header) map[string]string {
	headers := make(map[string]string)
	headers["x-request-id"] = header.Get("x-request-id")
	headers["x-b3-traceid"] = header.Get("x-b3-traceid")
	headers["x-b3-spanid"] = header.Get("x-b3-spanid")
	headers["x-b3-parentspanid"] = header.Get("x-b3-parentspanid")
	headers["x-b3-sampled"] = header.Get("x-b3-sampled")
	headers["x-b3-flags"] = header.Get("x-b3-flags")
	headers["Authorization"] = header.Get("Authorization")
	return headers
}

func makeURL(connection connection) string {
	return fmt.Sprintf("http://%s:%s/%s", connection.Name, connection.Port, connection.Path)
}

func handleReq(url, method string, headers map[string]string, client httpClient, ch chan response) {
	req, err := http.NewRequest(method, url, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		ch <- response{StatusCode: http.StatusServiceUnavailable}
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- response{StatusCode: http.StatusServiceUnavailable}
		return
	}
	response := response{}
	err = json.Unmarshal(bytes, &response)
	ch <- response
}

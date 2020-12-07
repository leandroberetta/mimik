package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Endpoint struct {
	Name        string       `json:"name"`
	Path        string       `json:"path"`
	Method      string       `json:"method"`
	Connections []Connection `json:"connections"`
}

type Connection struct {
	Name   string `json:"name"`
	Port   string `json:"port"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type Service struct {
	Name      string
	Port      string
	Endpoints []Endpoint
}

type Response struct {
	Name             string     `json:"name"`
	Version          string     `json:"version"`
	Path             string     `json:"path"`
	StatusCode       int        `json:"statusCode"`
	UpstreamResponse []Response `json:"upstreamResponse"`
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func main() {
	service, _ := NewService(
		os.Getenv("MIMIK_SERVICE_NAME"),
		os.Getenv("MIMIK_SERVICE_PORT"),
		os.Getenv("MIMIK_ENDPOINTS_FILE"))
	client := &http.Client{}
	http.HandleFunc("/", endpointHandler(service, client))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", service.Port), nil))
}

func NewService(name, port, fileName string) (Service, error) {
	service := Service{Name: name, Port: port}
	err := loadEndpoints(fileName, &service.Endpoints)

	return service, err
}

func loadEndpoints(fileName string, endpoints *[]Endpoint) error {
	file, err := os.Open(fileName)
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, endpoints)

	return err
}

func endpointHandler(service Service, client HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := Response{Name: service.Name, Version: "v1", StatusCode: http.StatusNotFound}
		ch := make(chan Response)
		for _, endpoint := range service.Endpoints {
			response.Path = r.URL.Path
			if endpoint.Path == r.URL.Path {
				response.StatusCode = http.StatusOK

				upstreamResponse := make([]Response, len(endpoint.Connections))

				for _, connection := range endpoint.Connections {
					go handleReq(makeURL(connection), connection.Method, client, ch)
				}

				for i, _ := range endpoint.Connections {
					upstreamResponse[i] = <-ch
				}

				response.UpstreamResponse = upstreamResponse
			}
		}

		responseJSON, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

func makeURL(connection Connection) string {
	return fmt.Sprintf("http://%s:%s/%s", connection.Name, connection.Port, connection.Path)
}

func handleReq(url, method string, client HTTPClient, ch chan Response) {
	req, err := http.NewRequest(method, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		ch <- Response{StatusCode: http.StatusServiceUnavailable}
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- Response{StatusCode: http.StatusServiceUnavailable}
		return
	}

	response := Response{}
	err = json.Unmarshal(bytes, &response)

	ch <- response
}

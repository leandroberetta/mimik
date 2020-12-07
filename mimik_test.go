package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMakeValidURL(t *testing.T) {
	connection := connection{Name: "songs-service", Port: "8080", Path: "songs/1"}
	expected := "http://songs-service:8080/songs/1"

	if got := makeURL(connection); got != expected {
		t.Errorf("Wrong URL: Expected %s - Got %s", expected, got)
	}
}

func TestGetVersion(t *testing.T) {
	expected := "v2"

	if got := getVersion("test/mimik_labels.txt"); got != expected {
		t.Errorf("Wrong version: Expected %s - Got %s", expected, got)
	}
}

func TestEndpointHandler(t *testing.T) {
	service, _ := newService("lyrics", "8080", "test/mimik_test.json", "test/mimik_labels.txt")

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	mockClient := NewMockHTTPClient(ctrl)

	response := response{Name: "songs-service", Version: "v1", Path: "/songs/1", StatusCode: 200, UpstreamResponse: []response{}}
	responseJSON, _ := json.Marshal(response)
	responseBytes := ioutil.NopCloser(bytes.NewReader(responseJSON))

	mockClient.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: responseBytes}, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(endpointHandler(service, mockClient))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: Expected: %v - Got: %v", http.StatusOK, status)
	}
}

func TestNewService(t *testing.T) {
	expectedServiceName := "lyrics"
	expectedServicePort := "8080"

	service, _ := newService(expectedServiceName, expectedServicePort, "test/mimik_test.json", "test/mimik_labels.txt")

	if gotServiceName := service.Name; gotServiceName != expectedServiceName {
		t.Errorf("Wrong service name: Expected %s - Got %s", expectedServiceName, gotServiceName)
	}

	if gotServicePort := service.Port; gotServicePort != expectedServicePort {
		t.Errorf("Wrong service port: Expected %s - Got %s", expectedServicePort, gotServicePort)
	}

	for _, endpoint := range service.Endpoints {
		if endpoint.Name == "/songs/1" {
			expectedConnections := 2

			if gotConnections := len(endpoint.Connections); gotConnections != expectedConnections {
				t.Errorf("Wrong number of connections: Expected %d - Got %d", expectedConnections, gotConnections)
			}

			for _, connection := range endpoint.Connections {
				if connection.Name == "songs-service" {
					expectedPort := "8080"
					expectedPath := "songs/1"
					expectedMethod := "GET"

					if gotPort := connection.Port; gotPort != expectedPort {
						t.Errorf("Wrong connection port: Expected %s - Got %s", expectedPort, gotPort)
					}
					if gotPath := connection.Path; gotPath != expectedPath {
						t.Errorf("Wrong connection path: Expected %s - Got %s", expectedPath, gotPath)
					}
					if gotMethod := connection.Method; gotMethod != expectedMethod {
						t.Errorf("Wrong connection method: Expected %s - Got %s", expectedMethod, gotMethod)
					}
				}
			}
		}
	}
}

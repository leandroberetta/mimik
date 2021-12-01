package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/leandroberetta/mimik/pkg/api"
	"github.com/leandroberetta/mimik/pkg/service"
)

func main() {
	log.Println("starting Mimik instance...")

	instance, _ := service.NewService(
		os.Getenv("MIMIK_SERVICE_NAME"),
		os.Getenv("MIMIK_SERVICE_PORT"),
		os.Getenv("MIMIK_ENDPOINTS_FILE"),
		service.GetVersion(os.Getenv("MIMIK_LABELS_FILE")))

	client := &http.Client{}

	r := mux.NewRouter()
	r.Path("/").HandlerFunc(service.EndpointHandler(instance, client))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("serving at :8080")
	go srv.ListenAndServe()

	tc := make(chan struct{})
	if tg := os.Getenv("MIMIK_TRAFFIC_GENERATOR"); tg != "" {
		go generateTraffic(&instance, client, tc)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	close(tc)

	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("shutdown complete")
	os.Exit(0)
}

func generateTraffic(service *api.Service, client *http.Client, quit chan struct{}) {
	for {
		select {
		case <-quit:
			log.Println("stopping traffic generator")
			return
		default:
			for _, endpoint := range service.Endpoints {
				req, _ := http.NewRequest(endpoint.Method, fmt.Sprintf("http://localhost:%d", 8080), nil)
				client.Do(req)
			}

			time.Sleep(1 * time.Second)
		}
	}
}

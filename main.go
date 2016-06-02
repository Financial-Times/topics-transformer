package main

import (
	"fmt"
	digest "github.com/FeNoMeNa/goha"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFormatter(new(log.JSONFormatter))
}

func main() {
	app := cli.App("topics-transformer", "A RESTful API for transforming TME Topics to UP json")
	username := app.String(cli.StringOpt{
		Name:   "structure-service-username",
		Value:  "",
		Desc:   "Structure service username used for http digest authentication",
		EnvVar: "STRUCTURE_SERVICE_USERNAME",
	})
	password := app.String(cli.StringOpt{
		Name:   "structure-service-password",
		Value:  "",
		Desc:   "Structure service password used for http digest authentication",
		EnvVar: "STRUCTURE_SERVICE_PASSWORD",
	})
	principalHeader := app.String(cli.StringOpt{
		Name:   "principal-header",
		Value:  "",
		Desc:   "Structure service principal header used for authentication",
		EnvVar: "PRINCIPAL_HEADER",
	})
	baseURL := app.String(cli.StringOpt{
		Name:   "base-url",
		Value:  "http://localhost:8080/transformers/topics/",
		Desc:   "Base url",
		EnvVar: "BASE_URL",
	})
	structureServiceBaseURL := app.String(cli.StringOpt{
		Name:   "structure-service-base-url",
		Value:  "http://metadata.internal.ft.com:83",
		Desc:   "Structure service base url",
		EnvVar: "STRUCTURE_SERVICE_BASE_URL",
	})
	port := app.Int(cli.IntOpt{
		Name:   "port",
		Value:  8080,
		Desc:   "Port to listen on",
		EnvVar: "PORT",
	})

	app.Action = func() {
		c := digest.NewClient(*username, *password)
		c.Timeout(10 * time.Second)
		s, err := newTopicService(newTmeRepository(c, *structureServiceBaseURL, *principalHeader), *baseURL)
		if err != nil {
			log.Errorf("Error while creating TopicsService: [%v]", err.Error())
		}
		h := newTopicsHandler(s)
		m := mux.NewRouter()

		// The top one of these feels more correct, but the lower one matches what we have in Dropwizard,
		// so it's what apps expect currently same as ping, the content of build-info needs more definition
		//using http router here to be able to catch "/"
		http.HandleFunc(status.PingPath, status.PingHandler)
		http.HandleFunc(status.PingPathDW, status.PingHandler)
		http.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
		http.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)
		http.HandleFunc("/__gtg", h.GoodToGo)

		m.HandleFunc("/transformers/topics", h.getTopics).Methods("GET")
		m.HandleFunc("/transformers/topics/{uuid}", h.getTopicByUUID).Methods("GET")

		http.Handle("/", m)

		log.Printf("listening on %d", *port)
		http.ListenAndServe(fmt.Sprintf(":%d", *port),
			httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry,
				httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), m)))
	}
	app.Run(os.Args)
}

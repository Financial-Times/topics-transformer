package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type topicsHandler struct {
	service topicService
}

// HealthCheck does something
func (h *topicsHandler) HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for the topic data from TME",
		Name:             "Check connectivity to TME",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/topics-transfomer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to TME to be able to supply topics",
		Checker:          h.checker,
	}
}

// Checker does more stuff
func (h *topicsHandler) checker() (string, error) {
	err := h.service.checkConnectivity()
	if err == nil {
		return "Connectivity to TME is ok", err
	}
	return "Error connecting to TME", err
}

func newTopicsHandler(service topicService) topicsHandler {
	return topicsHandler{service: service}
}

func (h *topicsHandler) getTopics(writer http.ResponseWriter, req *http.Request) {
	obj, found := h.service.getTopics()
	writeJSONResponse(obj, found, writer)
}

func (h *topicsHandler) getTopicByUUID(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	obj, found := h.service.getTopicByUUID(uuid)
	writeJSONResponse(obj, found, writer)
}

//GoodToGo returns a 503 if the healthcheck fails - suitable for use from varnish to check availability of a node
func (h *topicsHandler) GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := h.checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func writeJSONResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v\n", err)
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", errorMsg))
}

func (h *topicsHandler) getCount(writer http.ResponseWriter, req *http.Request) {
	count := h.service.getTopicCount()
	_, err := writer.Write([]byte(strconv.Itoa(count)))
	if err != nil {
		log.Warnf("Couldn't write count to HTTP response. count=%d %v\n", count, err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *topicsHandler) getIds(writer http.ResponseWriter, req *http.Request) {
	ids := h.service.getTopicIds()
	writer.Header().Add("Content-Type", "text/plain")
	if len(ids) == 0 {
		writer.WriteHeader(http.StatusOK)
		return
	}
	enc := json.NewEncoder(writer)
	type topicID struct {
		ID string `json:"id"`
	}
	for _, id := range ids {
		rID := topicID{ID: id}
		err := enc.Encode(rID)
		if err != nil {
			log.Warnf("Couldn't encode to HTTP response topic with uuid=%s %v\n", id, err)
			continue
		}
	}
}

func (h *topicsHandler) reload(writer http.ResponseWriter, req *http.Request) {
	err := h.service.reload()
	if err != nil {
		log.Warnf("Problem reloading terms from TME: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

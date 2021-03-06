package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getTopicsResponse = `[{"apiUrl":"http://localhost:8080/transformers/topics/bba39990-c78d-3629-ae83-808c333c6dbc"}]`
const getTopicByUUIDResponse = `{"uuid":"bba39990-c78d-3629-ae83-808c333c6dbc","alternativeIdentifiers":{"TME":["MTE3-U3ViamVjdHM="],"uuids":["bba39990-c78d-3629-ae83-808c333c6dbc"]},"prefLabel":"Metals Markets","type":"Topic"}`
const getTopicsCountResponse = `1`
const getTopicsIdsResponse = `{"id":"bba39990-c78d-3629-ae83-808c333c6dbc"}`

func TestHandlers(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		req          *http.Request
		dummyService topicService
		statusCode   int
		contentType  string // Contents of the Content-Type header
		body         string
	}{
		{"Success - get topic by uuid", newRequest("GET", fmt.Sprintf("/transformers/topics/%s", testUUID)), &dummyService{found: true, topics: []topic{getDummyTopic(testUUID, "Metals Markets", "MTE3-U3ViamVjdHM=")}}, http.StatusOK, "application/json", getTopicByUUIDResponse},
		{"Not found - get topic by uuid", newRequest("GET", fmt.Sprintf("/transformers/topics/%s", testUUID)), &dummyService{found: false, topics: []topic{topic{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get topics", newRequest("GET", "/transformers/topics"), &dummyService{found: true, topics: []topic{topic{UUID: testUUID}}}, http.StatusOK, "application/json", getTopicsResponse},
		{"Not found - get topics", newRequest("GET", "/transformers/topics"), &dummyService{found: false, topics: []topic{}}, http.StatusNotFound, "application/json", ""},
		{"Test Topic Count", newRequest("GET", "/transformers/topics/__count"), &dummyService{found: true, topics: []topic{topic{UUID: testUUID}}}, http.StatusOK, "text/plain", getTopicsCountResponse},
		{"Test Topic Ids", newRequest("GET", "/transformers/topics/__ids"), &dummyService{found: true, topics: []topic{topic{UUID: testUUID}}}, http.StatusOK, "text/plain", getTopicsIdsResponse},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.True(test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.Equal(strings.TrimSpace(test.body), strings.TrimSpace(rec.Body.String()), fmt.Sprintf("%s: Wrong body", test.name))
	}
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func router(s topicService) *mux.Router {
	m := mux.NewRouter()
	h := newTopicsHandler(s)
	m.HandleFunc("/transformers/topics", h.getTopics).Methods("GET")
	m.HandleFunc("/transformers/topics/__count", h.getCount).Methods("GET")
	m.HandleFunc("/transformers/topics/__ids", h.getIds).Methods("GET")
	m.HandleFunc("/transformers/topics/__reload", h.getIds).Methods("GET")
	m.HandleFunc("/transformers/topics/{uuid}", h.getTopicByUUID).Methods("GET")

	return m
}

type dummyService struct {
	found  bool
	topics []topic
}

func (s *dummyService) getTopics() ([]topicLink, bool) {
	var topicLinks []topicLink
	for _, sub := range s.topics {
		topicLinks = append(topicLinks, topicLink{APIURL: "http://localhost:8080/transformers/topics/" + sub.UUID})
	}
	return topicLinks, s.found
}

func (s *dummyService) getTopicByUUID(uuid string) (topic, bool) {
	return s.topics[0], s.found
}

func (s *dummyService) checkConnectivity() error {
	return nil
}

func (s *dummyService) getTopicCount() int {
	return len(s.topics)
}

func (s *dummyService) getTopicIds() []string {
	i := 0
	keys := make([]string, len(s.topics))

	for _, t := range s.topics {
		keys[i] = t.UUID
		i++
	}
	return keys
}

func (s *dummyService) reload() error {
	return nil
}

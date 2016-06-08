package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getTopicsResponse = `[{"apiUrl":"http://localhost:8080/transformers/topics/bba39990-c78d-3629-ae83-808c333c6dbc"}]`
const getTopicByUUIDResponse = `{"uuid":"bba39990-c78d-3629-ae83-808c333c6dbc","alternativeIdentifiers":{"TME":["MTE3-U3ViamVjdHM="],"uuids":["bba39990-c78d-3629-ae83-808c333c6dbc"]},"prefLabel":"Metals Markets","type":"Topic"}`

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

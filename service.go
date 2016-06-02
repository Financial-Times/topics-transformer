package main

import (
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type topicService interface {
	getTopics() ([]topicLink, bool)
	getTopicByUUID(uuid string) (topic, bool)
	checkConnectivity() error
}

type topicServiceImpl struct {
	repository repository
	baseURL    string
	topicsMap  map[string]topic
	topicLinks []topicLink
}

func newTopicService(repo repository, baseURL string) (topicService, error) {

	s := &topicServiceImpl{repository: repo, baseURL: baseURL}
	err := s.init()
	if err != nil {
		return &topicServiceImpl{}, err
	}
	return s, nil
}

func (s *topicServiceImpl) init() error {
	s.topicsMap = make(map[string]topic)
	tax, err := s.repository.getTopicsTaxonomy()
	if err != nil {
		return err
	}
	s.initTopicsMap(tax.Terms)
	return nil
}

func (s *topicServiceImpl) getTopics() ([]topicLink, bool) {
	if len(s.topicLinks) > 0 {
		return s.topicLinks, true
	}
	return s.topicLinks, false
}

func (s *topicServiceImpl) getTopicByUUID(uuid string) (topic, bool) {
	topic, found := s.topicsMap[uuid]
	return topic, found
}

func (s *topicServiceImpl) checkConnectivity() error {
	// TODO: Can we just hit an endpoint to check if TME is available? Or do we need to make sure we get genre taxonmies back?
	_, err := s.repository.getTopicsTaxonomy()
	if err != nil {
		return err
	}
	return nil
}

func (s *topicServiceImpl) initTopicsMap(terms []term) {
	for _, t := range terms {
		top := transformTopic(t)
		s.topicsMap[top.UUID] = top
		s.topicLinks = append(s.topicLinks, topicLink{APIURL: s.baseURL + top.UUID})
		s.initTopicsMap(t.Children.Terms)
	}
}

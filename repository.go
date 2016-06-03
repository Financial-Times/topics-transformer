package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type repository interface {
	getTopicsTaxonomy() (taxonomy, error)
}

type tmeRepository struct {
	httpClient              httpClient
	principalHeader         string
	structureServiceBaseURL string
}

func newTmeRepository(client httpClient, structureServiceBaseURL string, principalHeader string) repository {
	return &tmeRepository{httpClient: client, principalHeader: principalHeader, structureServiceBaseURL: structureServiceBaseURL}
}

func (t *tmeRepository) getTopicsTaxonomy() (taxonomy, error) {
	requestString := t.structureServiceBaseURL + "/metadata-services/structure/1.0/taxonomies/topics/terms?includeDisabledTerms=false"
	req, err := http.NewRequest("GET", requestString, nil)
	if err != nil {
		return taxonomy{}, err
	}
	req.Header.Set("ClientUserPrincipal", t.principalHeader)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return taxonomy{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return taxonomy{}, fmt.Errorf("Structure service returned a %d http status code while connecting to %s", resp.StatusCode, requestString)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return taxonomy{}, err
	}

	tax := taxonomy{}
	err = xml.Unmarshal(contents, &tax)
	if err != nil {
		return taxonomy{}, err
	}
	return tax, nil
}

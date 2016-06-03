package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestGetTopicsTaxonomy(t *testing.T) {
	assert := assert.New(t)
	topicsXML, _ := os.Open("sample_topics.xml")
	tests := []struct {
		name string
		repo repository
		tax  taxonomy
		err  error
	}{
		{"Success", repo(dummyClient{assert: assert, structureServiceBaseURL: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(topicsXML)}}),
			taxonomy{Terms: []term{term{CanonicalName: "Z_Archive", ID: "NWQxOTE4ZGQtOGY1OS00MWY3LTk0ZWEtOWYyOGNmMDg4ZGJk-VG9waWNz", Children: children{[]term{term{CanonicalName: "Africa Inc", ID: "YTcyNWI5YzItOTUwMy00ZWRkLWI0M2YtYzBjZjU5MWNjNTJi-VG9waWNz"}}}}}}, nil},
		{"Error", repo(dummyClient{assert: assert, structureServiceBaseURL: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(topicsXML)}, err: errors.New("Some error")}),
			taxonomy{}, errors.New("Some error")},
		{"Non 200 from structure service", repo(dummyClient{assert: assert, structureServiceBaseURL: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(topicsXML)}}),
			taxonomy{}, errors.New("Structure service returned a 400 http status code while connecting to http://metadata.internal.ft.com:83/metadata-services/structure/1.0/taxonomies/topics/terms?includeDisabledTerms=false")},
		{"Unmarshalling error", repo(dummyClient{assert: assert, structureServiceBaseURL: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte("Non xml")))}}),
			taxonomy{}, errors.New("EOF")},
	}

	for _, test := range tests {
		expectedTax, err := test.repo.getTopicsTaxonomy()
		assert.Equal(test.tax, expectedTax, fmt.Sprintf("%s: Expected taxonomy incorrect", test.name))
		assert.Equal(test.err, err)
	}

}

func repo(c dummyClient) repository {
	return &tmeRepository{httpClient: &c, principalHeader: c.principalHeader, structureServiceBaseURL: c.structureServiceBaseURL}
}

type dummyClient struct {
	assert                  *assert.Assertions
	resp                    http.Response
	err                     error
	principalHeader         string
	structureServiceBaseURL string
}

func (d *dummyClient) Do(req *http.Request) (resp *http.Response, err error) {
	d.assert.Equal(d.principalHeader, req.Header.Get("ClientUserPrincipal"), fmt.Sprintf("Expected ClientUserPrincipal header incorrect"))
	d.assert.Equal(d.structureServiceBaseURL+"/metadata-services/structure/1.0/taxonomies/topics/terms?includeDisabledTerms=false", req.URL.String(), fmt.Sprintf("Expected url incorrect"))
	return &d.resp, d.err
}

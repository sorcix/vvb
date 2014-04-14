package vvb

import (
	"io"
	"net/http"
)

var Client = http.DefaultClient

func get(url string) (body io.ReadCloser, err error) {

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	request.Header.Add("Accept", "application/xml,text/xml")
	request.Header.Add("User-Agent", "VolleyVVB-Go (+http://github.com/sorcix/vvb)")
	request.Header.Add("Content-Type", "application/xml")

	response, err := Client.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode == 200 {
		return response.Body, nil
	}

	return

}

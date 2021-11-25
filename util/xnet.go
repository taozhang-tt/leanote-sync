package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func HttpRequest(address string, param map[string]string) ([]byte, error) {
	var (
		body      = []byte{}
		urlValues = url.Values{}
	)
	for k, v := range param {
		urlValues.Set(k, v)
	}
	resp, err := http.PostForm(address, urlValues)
	if err != nil {
		return body, err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("request %s err: %s", address, resp.Status)
		return body, errors.New("HttpRequest Response Code: " + resp.Status)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var HttpClient *http.Client

func init() {
	HttpClient = http.DefaultClient
}

func PostRawJson(finalURL string, req []byte, response interface{}) (err error) {
	httpResp, err := HttpClient.Post(finalURL, "application/json; charset=utf-8", bytes.NewReader(req))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	if err = json.NewDecoder(httpResp.Body).Decode(response); err != nil {
		return
	}
	return
}

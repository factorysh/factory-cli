package client

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func readJson(resp *http.Response, value interface{}) error {
	if resp.StatusCode%100 != 2 {
		return fmt.Errorf("Bad status: %v", resp.Status)
	}
	var reader io.ReadCloser
	enc := resp.Header.Get("Content-encoding")
	if enc != "" {
		if enc == "gzip" {
			var err error
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Wrong encoding: %s", enc)
		}
	} else {
		reader = resp.Body
	}
	defer reader.Close()
	raw, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, value)
}

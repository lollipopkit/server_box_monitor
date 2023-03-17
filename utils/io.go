package utils

import (
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	httpClient = http.DefaultClient
)

func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func HttpDo(method, url, content string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(content))
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}

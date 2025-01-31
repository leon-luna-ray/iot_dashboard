package utils

import (
	"io"
	"net/http"
)

func CopyResponse(w http.ResponseWriter, resp *http.Response) error {
	defer resp.Body.Close()
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	_, err := io.Copy(w, resp.Body)
	return err
}

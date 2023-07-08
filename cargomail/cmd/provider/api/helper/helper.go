package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func ReturnErr(w http.ResponseWriter, err error, code int) {
	errorMessage := struct {
		Err string
	}{
		Err: err.Error(),
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorMessage)
}

func SetJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func FromJson[T any](body io.Reader, target T) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	json.Unmarshal(buf.Bytes(), &target)
}


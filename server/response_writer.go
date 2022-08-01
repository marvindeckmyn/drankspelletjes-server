package server

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter struct {
	W http.ResponseWriter
}

func newRespWriter(w http.ResponseWriter) ResponseWriter {
	return ResponseWriter{
		W: w,
	}
}

func (rw *ResponseWriter) JSON(status int, content interface{}) error {
	rw.W.Header().Set("Content-Type", "application/json")

	if content != nil {
		data, err := json.Marshal(content)
		if err != nil {
			rw.W.Write(data)
			return &ErrMarshaling{}
		}

		rw.W.WriteHeader(status)
		rw.W.Write(data)
	} else {
		rw.W.WriteHeader(status)
		rw.W.Write([]byte("{}"))
	}

	return nil
}

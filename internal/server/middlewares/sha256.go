package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

type sha256ResponseWriter struct {
	w    http.ResponseWriter
	hash string
}

func (sw *sha256ResponseWriter) Header() http.Header {
	return sw.w.Header()
}
func (sw *sha256ResponseWriter) WriteHeader(code int) {
	sw.w.WriteHeader(code)
}
func (sw *sha256ResponseWriter) Write(data []byte) (int, error) {
	sw.w.Header().Set("HashSHA256", sw.hash)

	return sw.w.Write(data)
}

func WithSha256(hf http.HandlerFunc, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf.Bytes()))

		h := hmac.New(sha256.New, []byte(key))
		if _, err := h.Write(buf.Bytes()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash := h.Sum(nil)

		hashStr := hex.EncodeToString(hash)
		headerHashStr, isHashHeaderSet := r.Header["HashSHA256"]
		if isHashHeaderSet && hashStr != headerHashStr[0] {
			http.Error(w, "Hash is invalid", http.StatusNotImplemented)
			return
		}

		hf(&sha256ResponseWriter{
			w:    w,
			hash: hashStr,
		}, r)
	}
}

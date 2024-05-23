package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
)

type sha256ResponseWriter struct {
	w    http.ResponseWriter
	hash []byte
}

func (sw *sha256ResponseWriter) Header() http.Header {
	return sw.w.Header()
}
func (sw *sha256ResponseWriter) WriteHeader(code int) {
	sw.w.WriteHeader(code)
}
func (sw *sha256ResponseWriter) Write(data []byte) (int, error) {
	sw.w.Header().Set("HashSHA256", string(sw.hash))

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

		h := hmac.New(sha256.New, []byte(key))
		if _, err := h.Write(buf.Bytes()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash := h.Sum(nil)

		if string(hash) != r.Header.Get("HashSHA256") {
			http.Error(w, "Hash is invalid", http.StatusBadRequest)
			return
		}

		hf(&sha256ResponseWriter{
			w:    w,
			hash: hash,
		}, r)
	}
}

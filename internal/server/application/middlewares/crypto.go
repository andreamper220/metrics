package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"os"
)

func WithCrypto(hf http.HandlerFunc, keyPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		privateKeyPEM, err := os.ReadFile(keyPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		privateKeyBlock, _ := pem.Decode(privateKeyPEM)
		privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, buf.Bytes())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(plaintext))

		hf(w, r)
	}
}

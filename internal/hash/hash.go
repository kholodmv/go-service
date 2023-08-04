package hash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

func HashHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		receivedHash := fmt.Sprintf("%x", sha256.Sum256(body))
		headerHash := r.Header.Get("HashSHA256")

		if headerHash != receivedHash {
			http.Error(w, "Incorrect HashSHA256 header value", http.StatusBadRequest)
			return
		}

		response := string(body)

		calculatedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(response)))
		w.Header().Set("HashSHA256", calculatedHash)

		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

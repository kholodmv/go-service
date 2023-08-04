package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

func HashHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hashing")
		headerHashValue := r.Header.Get("HashSHA256")
		if headerHashValue == "" {
			next.ServeHTTP(w, r)
			return
		}
		body, _ := io.ReadAll(r.Body)
		fmt.Println(string(body))
		defer r.Body.Close()

		h := hmac.New(sha256.New, []byte("kek"))
		h.Write(body)
		expectedHash := h.Sum(nil)

		expectedHashString := fmt.Sprintf("%x", expectedHash)
		fmt.Println(expectedHashString)
		fmt.Println(headerHashValue)
		if expectedHashString == headerHashValue {
			fmt.Println("hashing ServeHTTP")
			r.Body = io.NopCloser(bytes.NewReader(body))
			fmt.Println(r.Body)
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Incorrect HashSHA256 header value"))
		}
		next.ServeHTTP(w, r)
	})

	/*body, err := io.ReadAll(r.Body)
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
	next.ServeHTTP(w, r)*/
}

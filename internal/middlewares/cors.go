package middlewares

import (
	"net/http"
)

func CorsMiddleware(ah http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")                                              // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")               // Allow specific methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With") // Allow headers
		w.Header().Set("Access-Control-Allow-Credentials", "true")                                      // Allow cookies (if needed)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		ah.ServeHTTP(w, r)
	})
}

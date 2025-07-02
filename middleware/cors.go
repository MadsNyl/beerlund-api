package middleware

import "net/http"

// corsMiddleware injects the CORS headers and handles preflight.
func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1) Allow your frontend origin (or use "*" if you don't care)
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        // 2) Allow the methods you support
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        // 3) Allow any headers your client will send
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        // 4) (Optional) expose any headers back to the client
        w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

        // If this is a preflight OPTIONS request, we stop here:
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        // Otherwise, call the next handler in chain:
        next.ServeHTTP(w, r)
    })
}
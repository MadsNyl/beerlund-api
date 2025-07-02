package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

func AuthMiddleware(next http.Handler, clerkClient *user.Client) http.Handler {
    // Build and apply the header-auth middleware
    mw := clerkhttp.WithHeaderAuthorization()
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1) verify the token + populate claims in ctx
        base := mw(next)
        base.ServeHTTP(w, r)

        // 2) extract the SessionClaims from the *root* clerk package
        claims, ok := clerk.SessionClaimsFromContext(r.Context())
        log.Printf("Session claims: %+v\n", claims)
    
        if !ok {
            http.Error(w, "Unauthorized – no valid session", http.StatusUnauthorized)
            return
        }

        userID := claims.Subject
        if userID == "" {
            http.Error(w, "Unauthorized – missing user ID", http.StatusUnauthorized)
            return
        }

        // 3) optional: verify user exists
        if _, err := clerkClient.Get(r.Context(), userID); err != nil {
            http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusUnauthorized)
            return
        }

        // 4) inject into context for downstream handlers
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func ProtectedRoute(w http.ResponseWriter, r *http.Request) {
  claims, ok := clerk.SessionClaimsFromContext(r.Context())
  if !ok {
    w.WriteHeader(http.StatusUnauthorized)
    w.Write([]byte(`{"access": "unauthorized"}`))
    return
  }

  usr, err := user.Get(r.Context(), claims.Subject)
  if err != nil {
    // handle the error
  }
  fmt.Fprintf(w, `{"user_id": "%s", "user_banned": "%t"}`, usr.ID, usr.Banned)
}


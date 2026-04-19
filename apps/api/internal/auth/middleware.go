// Package auth provides Auth0 JWT validation middleware for Chi.
package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type contextKey string

const claimsKey contextKey = "claims"

// Claims holds the decoded JWT payload fields we care about.
type Claims struct {
	Sub   string // Auth0 subject (user ID)
	Email string // email claim (optional, requires profile scope)
}

// Middleware validates Auth0 JWTs on every request.
// Unauthenticated requests receive 401; the decoded claims are stored in ctx.
type Middleware struct {
	jwksURL  string
	audience string
	issuer   string
	cache    *jwk.Cache
}

// New creates a Middleware for the given Auth0 domain and API audience.
func New(domain, audience string) *Middleware {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)
	cache := jwk.NewCache(context.Background())
	// Register the JWKS URL for automatic refresh.
	if err := cache.Register(jwksURL); err != nil {
		slog.Error("failed to register JWKS URL", "url", jwksURL, "err", err)
	}
	return &Middleware{
		jwksURL:  jwksURL,
		audience: audience,
		issuer:   fmt.Sprintf("https://%s/", domain),
		cache:    cache,
	}
}

// Handler returns an http.Handler middleware that validates Bearer tokens.
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := bearerToken(r)
		if !ok {
			http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
			return
		}

		keySet, err := m.cache.Get(r.Context(), m.jwksURL)
		if err != nil {
			slog.Error("failed to fetch JWKS", "err", err)
			http.Error(w, "could not validate token", http.StatusInternalServerError)
			return
		}

		token, err := jwt.Parse([]byte(raw),
			jwt.WithKeySet(keySet),
			jwt.WithValidate(true),
			jwt.WithAudience(m.audience),
			jwt.WithIssuer(m.issuer),
		)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims := Claims{Sub: token.Subject()}
		if email, ok := token.Get("email"); ok {
			claims.Email, _ = email.(string)
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ClaimsFromContext retrieves the validated claims from a request context.
// Returns zero value and false if not present (i.e. unauthenticated route).
func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsKey).(Claims)
	return c, ok
}

func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", false
	}
	return parts[1], true
}

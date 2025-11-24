package main

import (
	"net/http"

	"github.com/ctiller15/tailscribe/internal/auth"
)

type baseHandler func(w http.ResponseWriter, r *http.Request)

func (a *APIConfig) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		a.Logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next(w, r)
	})
}

type authorizedHandler func(w http.ResponseWriter, r *http.Request, userID int)

func (a *APIConfig) CheckAuthMiddleware(handler authorizedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check user's cookie
		jwtCookie := r.CookiesNamed("token")

		if len(jwtCookie) == 0 {
			a.Logger.Error("invalid request, user unauthorized, no token.")
			// TODO: attempt to refresh before redirecting and potentially create a new token.
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		tokenString := jwtCookie[0].Value
		user_id, err := auth.ValidateJWT(tokenString, a.Env.Secret)
		if err != nil {
			a.Logger.Error("invalid request, user token invalid")
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		handler(w, r, user_id)
	}
}

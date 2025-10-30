package api

import (
	"log"
	"net/http"

	"github.com/ctiller15/tailscribe/internal/auth"
)

type authorizedHandler func(w http.ResponseWriter, r *http.Request, userID int)

func (a *APIConfig) CheckAuthMiddleware(handler authorizedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check user's cookie
		jwtCookie := r.CookiesNamed("token")

		if len(jwtCookie) == 0 {
			log.Printf("invalid request, user unauthorized, no token.")
			// TODO: attempt to refresh before redirecting and potentially create a new token.
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		tokenString := jwtCookie[0].Value
		user_id, err := auth.ValidateJWT(tokenString, a.Env.Secret)
		if err != nil {
			log.Printf("invalid request, user token invalid")
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		handler(w, r, user_id)
	}
}

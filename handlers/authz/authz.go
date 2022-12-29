// Package authz for azure authorization checks
package authz

import (
	"context"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/zicops/sidecar-auth-proxy/lib/googleprojectlib"
	"github.com/zicops/sidecar-auth-proxy/lib/identity"
	"github.com/zicops/sidecar-auth-proxy/lib/jwt"
)

// Auth ...
var Auth *identity.IDP

// Check checks if the user is authenticated
func Check(h http.Handler) http.Handler {
	ctxAuth := context.Background()
	currentProject := googleprojectlib.GetGoogleProjectDefaultID()
	Auth, err := identity.NewIDPEP(ctxAuth, currentProject)
	if err != nil {
		log.Panicf("Auth: %s", err.Error())
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Authz in process")
		if strings.Contains(r.URL.Path, "healthz") {
			log.Errorf("Method does not exist. /healthz should not be visible from outside cluster")
			http.Error(w, "Method does not exist.", http.StatusUnauthorized)
			return
		}
		if strings.Contains(r.URL.Path, "reset-password") || strings.Contains(r.URL.Path, "org") {
			h.ServeHTTP(w, r)
			return
		}
		if Auth == nil {
			http.Error(w, "Fatal: Failed to initialize auth.", http.StatusInternalServerError)
			return
		}
		incomingToken := jwt.GetToken(r)
		if incomingToken == "" {
			incomingToken = jwt.GetTokenWebsocket(r)
		}
		claimsFromToken, err := jwt.GetClaims(incomingToken)
		if claimsFromToken == nil {
			log.Errorf("Failed to get claims from token: %s", err.Error())
			claimsFromToken = make(map[string]interface{})
		}
		ctx := context.WithValue(r.Context(), "zclaims", claimsFromToken)

		returnedToken, err := Auth.VerifyUserToken(ctx, incomingToken)
		if err != nil && returnedToken == nil {
			log.Errorf("Token signature verification failed. Error: %v", err)
			http.Error(w, "Unauthorized: Bad request or authorization details, invalid token", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

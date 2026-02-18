package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	request "github.com/golang-jwt/jwt/v4/request"
	"github.com/keel-hq/keel/pkg/auth"
	log "github.com/sirupsen/logrus"
)

func authHeadersMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	rw.Header().Set("Access-Control-Expose-Headers", "Authorization")
	rw.Header().Set("Access-Control-Request-Headers", "Authorization")

	next(rw, r)
}

func (s *TriggerServer) requireAdminAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		// rw.Header().Set("Access-Control-Expose-Headers", "Authorization")
		// rw.Header().Set("Access-Control-Request-Headers", "Authorization")
		//
		if r.Method == "OPTIONS" {
			rw.WriteHeader(200)
			return
		}

		username, password, ok := r.BasicAuth()
		if ok {
			resp, err := s.authenticator.Authenticate(&auth.AuthRequest{
				Username: username,
				Password: password,
				AuthType: auth.AuthTypeBasic,
			})

			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"user":  username,
					"pas":   password,
				}).Error("failed uath")
				// rw.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			r = auth.SetAuthenticationDetails(r, &resp.User)

			next(rw, r)
			return
		}

		// authenticating via token

		resp, err := s.authenticator.Authenticate(&auth.AuthRequest{
			Token:    extractToken(r),
			AuthType: auth.AuthTypeToken,
		})

		if err != nil {
			// rw.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

			log.Warnf("authentication by token failed, token: %s, err: %s", extractToken(r), err)
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		r = auth.SetAuthenticationDetails(r, &resp.User)

		next(rw, r)
	}
}

func extractToken(req *http.Request) string {
	ex := request.AuthorizationHeaderExtractor
	token, err := ex.ExtractToken(req)
	if err != nil {
		return ""
	}

	return token
}

// logoutHandler godoc
// @Summary Logout user
// @Description Logs out the currently authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {string} string "Unauthorized"
// @Security ApiKeyAuth
// @Router /v1/auth/logout [post]
// @Router /v1/auth/logout [get]
func (s *TriggerServer) logoutHandler(resp http.ResponseWriter, req *http.Request) {

	resp.WriteHeader(200)
	resp.Write([]byte(`{}`))
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// loginResponse represents the authentication response
type loginResponse struct {
	Token string `json:"token"`
}

// loginHandler godoc
// @Summary Login user
// @Description Authenticates a user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body loginRequest true "Login credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Username or password incorrect"
// @Router /v1/auth/login [post]
func (s *TriggerServer) loginHandler(resp http.ResponseWriter, req *http.Request) {

	var lr loginRequest
	dec := json.NewDecoder(req.Body)
	defer req.Body.Close()

	err := dec.Decode(&lr)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(resp, "%s", err)
		return
	}

	authResp, err := s.authenticator.Authenticate(&auth.AuthRequest{
		Username: lr.Username,
		Password: lr.Password,
		AuthType: auth.AuthTypeBasic,
	})

	if err != nil {
		log.Warnf("auth failed for user '%s', error: %s", lr.Username, err)
		http.Error(resp, "username or password incorrect", 401)
		return
	}

	log.Infof("auth successful for user %s", lr.Username)

	resp.Header().Add("Access-Control-Expose-Headers", "Authorization")
	resp.Header().Add("Authorization", fmt.Sprintf("Bearer %s", authResp.Token))

	response(authResp, 200, nil, resp, req)
}

// refreshHandler godoc
// @Summary Refresh authentication token
// @Description Generates a new authentication token for the current user
// @Tags auth
// @Produce json
// @Success 200 {object} loginResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /v1/auth/refresh [get]
func (s *TriggerServer) refreshHandler(resp http.ResponseWriter, req *http.Request) {
	user := auth.GetAccountFromCtx(req.Context())

	authResp, err := s.authenticator.GenerateToken(*user)
	if err != nil {
		response(nil, http.StatusOK, err, resp, req)
		return
	}

	// adding token to header
	resp.Header().Add("Access-Control-Expose-Headers", "Authorization")
	resp.Header().Add("Authorization", fmt.Sprintf("Bearer %s", authResp.Token))

	response(authResp, http.StatusOK, err, resp, req)
}

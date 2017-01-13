package reportportal

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

//User represents logged-in user
type User struct {
	Name        string
	Authorities []string
}

//UserInfoErr represents error response from ReportPortal's UAT endpoint
type UserInfoErr struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

//authError represents error (response and status code) from ReportPortal's UAT endpoint
type authError struct {
	errorDesc  *UserInfoErr
	statusCode int
}

//Error represents implementation of default golang's Error interface
func (e *authError) Error() string {
	r, err := json.Marshal(e.errorDesc)
	if nil != e {
		return err.Error()
	}
	return string(r)
}

const (
	authorizationHeader    = "Authorization"
	contentTypeHeader      = "Content-Type"
	jsonContentType        = "application/json"
	bearerToken            = "bearer"
	unknownAuthorityWeight = 0
)

//Authorities represents available ReportPortal roles
var Authorities = map[string]int{
	"ROLE_USER":          1,
	"ROLE_ADMINISTRATOR": 2,
}

//RequireRole checks whether request auth represented by ReportPortal user with provided or higher role
func RequireRole(role string, authServerURL string) gin.HandlerFunc {
	authority := "ROLE_" + strings.ToUpper(role)
	return func(c *gin.Context) {
		token, err := parseBearer(c.Request)

		if err != nil || token == "" {
			notAuthorized(c)
			return
		}

		info, err := getTokenInfo(token, authServerURL)
		if err != nil {
			authErr, ok := err.(*authError)
			if !ok {
				notAuthorized(c)
			} else {
				respondWithError(c, authErr.statusCode, authErr.errorDesc)
			}
			return
		}

		if !hasAuthority(authority, info.Authorities) {
			notAuthorized(c)
			return
		}

		c.Next()
	}
}

//notAuthorized sends 401 error to the client
func notAuthorized(c *gin.Context) {
	respondWithErrorString(c, http.StatusUnauthorized, "Not Authorized")
}

//respondWithErrorString wraps error with JSON ans sends 401 to the client
func respondWithErrorString(c *gin.Context, code int, message string) {
	respondWithError(c, code, gin.H{"error": message})
}

//respondWithErrorString converts message JSON ans sends 401 to the client
func respondWithError(c *gin.Context, code int, message interface{}) {
	c.JSON(code, message)
	c.Abort()
}

//parseBearer parses authorization header
func parseBearer(r *http.Request) (string, error) {
	authHeader := r.Header.Get(authorizationHeader)
	if "" != authHeader {
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || bearerToken != strings.ToLower(authHeaderParts[0]) {
			return "", fmt.Errorf("Authorization header format must be '%s: token'", bearerToken)
		}
		return authHeaderParts[1], nil

	}
	return r.URL.Query().Get("access_token"), nil
}

//getTokenInfo - obtains token info from ReportPortal's UAT service
func getTokenInfo(token string, authServerURL string) (*User, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	rq, _ := http.NewRequest("GET", authServerURL, nil)
	rq.Header.Add(authorizationHeader, fmt.Sprintf("%s %s", bearerToken, token))
	rq.Header.Add(contentTypeHeader, jsonContentType)
	rs, e := netClient.Do(rq)
	defer rs.Body.Close()

	if nil != e {
		return nil, e
	}

	if rs.StatusCode/100 > 2 {
		uatErr := new(UserInfoErr)
		decodeJSON(rs, uatErr)
		e = &authError{
			errorDesc:  uatErr,
			statusCode: rs.StatusCode,
		}
		return nil, e
	}
	user := new(User)
	decodeJSON(rs, user)
	return user, nil
}

func decodeJSON(rs *http.Response, v interface{}) error {
	return json.NewDecoder(rs.Body).Decode(v)
}

//hasAuthority checks whether user authorities has at least one which has equal or higher weight than expected authority
func hasAuthority(ea string, ua []string) bool {
	weight := Authorities[ea]
	//Role is unknown
	if unknownAuthorityWeight == weight {
		return false
	}

	//go through
	for _, r := range ua {
		if Authorities[r] >= weight {
			return true
		}
	}
	return false
}

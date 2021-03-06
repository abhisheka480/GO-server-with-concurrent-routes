package JwtAuthentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GO-server-with-concurrent-routes/config"
	"github.com/GO-server-with-concurrent-routes/models"
	"github.com/dgrijalva/jwt-go"
)

func JwtTokenSet(w http.ResponseWriter, r *http.Request) {
	creds := &models.Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Request Body could not be decoded")
		resData := createResponse(http.StatusBadRequest, "Request Body could not be decoded")
		json.NewEncoder(w).Encode(resData)
		return
	}
	fmt.Println("creds:", creds)

	expectedPassword, ok := config.Users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Wrong token sent")
		resData := createResponse(http.StatusUnauthorized, "Wrong password")
		json.NewEncoder(w).Encode(resData)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(config.MYSECRETKEY)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Wrong token sent")
		resData := createResponse(http.StatusBadRequest, "Wrong token sent")
		json.NewEncoder(w).Encode(resData)
		return
	}

	http.SetCookie(w, &http.Cookie{ //set the token as a cookie in the header of request
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Println("Token Successfully Set")
	resData := createResponse(http.StatusOK, "Token Successfully Set")
	json.NewEncoder(w).Encode(resData)
}

func AuthenticateUser(next http.Handler) http.Handler { //this will act as a middleware in incoming requests to validate user
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token") //get the jwt token set as cookie in the request header
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Println("Unauthorised")
				resData := createResponse(http.StatusUnauthorized, "Unauthorised")
				json.NewEncoder(w).Encode(resData)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Unauthorised")
			resData := createResponse(http.StatusUnauthorized, "Unauthorised")
			json.NewEncoder(w).Encode(resData)
			return
		}

		// Get the JWT string from the cookie
		tokenString := c.Value

		claims := &models.Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.MYSECRETKEY, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Println("Unauthorised")
				resData := createResponse(http.StatusUnauthorized, "Unauthorised")
				json.NewEncoder(w).Encode(resData)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Unauthorised")
			resData := createResponse(http.StatusUnauthorized, "Unauthorised")
			json.NewEncoder(w).Encode(resData)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Token is invalid")
			resData := createResponse(http.StatusUnauthorized, "Token is invalid")
			json.NewEncoder(w).Encode(resData)
			return
		}
		fmt.Println("User is authorised")
		next.ServeHTTP(w, r)
	})
}

func createResponse(code int, message string) models.ResponseData {
	responseData := models.ResponseData{}
	responseData.StatusCode = code
	responseData.Message = message
	return responseData
}

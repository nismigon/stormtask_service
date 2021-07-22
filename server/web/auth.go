package web

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

type Credentials struct {
	Email    string
	Password string
}

// GenerateToken generates a token with the provided credentials
// In the nominal case, this function generates a string corresponding to the token
// If the credentials are not found in the database, this returns an empty string
// If an error occurred, this returns the error
func (s *Server) generateToken(cred Credentials) (string, error) {
	user, err := s.Database.Authenticate(cred.Email, cred.Password)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil
	}
	// Generation of the claim object
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		IsAdmin: user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(s.Configuration.JWTSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateAndExtractToken get the token and extract it. It also verify if the token is valid.
// In the nominal case, this returns a Claim object with all the information
// If the token isn't valid, this returns nil
// If an error occurred, this returns the erro
func (s *Server) ValidateAndExtractToken(token string) (*Claims, error) {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Configuration.JWTSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	// Valid token verification
	if !parsedToken.Valid {
		return nil, nil
	}
	return claims, nil
}

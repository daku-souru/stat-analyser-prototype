package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	model "github.com/daku-souru/stat-analyser-api/model"
)

func ProcessVerdictRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract token from header
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if Token is signed
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid Signing algorithm")
		}

		// Secret key for checking the token
		return []byte(model.SHARED_SECRET), nil
	})

	// Check if valid
	if err != nil {
		log.Print(err)
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	// Get the claims from the JWT token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Print(err)
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	// Check if the token has expired
	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		log.Print(err)
		http.Error(w, "Token has expired", http.StatusUnauthorized)
		return
	}

	if claims["iss"] != "Client" {
		log.Print(err)
		http.Error(w, "Wrong issuer", http.StatusUnauthorized)
		return
	}

	if claims["sub"] != "Match" {
		log.Print(err)
		http.Error(w, "Wrong subject", http.StatusUnauthorized)
		return
	}

	var verdict model.Verdict
	err = json.NewDecoder(r.Body).Decode(&verdict)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	if verdict.AutomationResult > model.V_OK || verdict.StatResult > model.V_OK || verdict.AutomationResult == model.V_CORRUPTED || verdict.StatResult == model.V_CORRUPTED {
		err = model.InsertVerdict(verdict)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Verdict was successfully added")
}

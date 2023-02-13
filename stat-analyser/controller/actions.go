package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/daku-souru/stat-analyser/model"
)

const API = "http://localhost:8080/stat-analyser-api"

func ConvertVerdictToJSON(v model.Verdict) string {
	//json.MarshalIndent(v, "", "\t")
	j, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(j)
}

func SendVerdict(verdict model.Verdict) error {
	secret := []byte(model.SHARED_SECRET)

	// Claims to be encoded in the JWT token
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Issuer:    "Client",
		Subject:   "Match",
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the JWT token with the secret key
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return err
	}

	jsonData := ConvertVerdictToJSON(verdict)

	// Create a new HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Create a new POST request with the JSON data
	req, err := http.NewRequest("POST", API, bytes.NewBufferString(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", signedToken)

	// Send the POST request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Body)
		return fmt.Errorf("Error sending request: %s", resp.Status)
	}

	return nil
}

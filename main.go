package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/resty.v1"
)

// CLI struct for environment variables and CLI arguments
type CLI struct {
	KeycloakURL  string `help:"Keycloak server URL" env:"KEYCLOAK_URL" default:"http://localhost:8080"`
	Realm        string `help:"Keycloak realm" env:"KEYCLOAK_REALM" default:"master"`
	ClientID     string `help:"Keycloak client ID" env:"KEYCLOAK_CLIENT_ID" default:"user-manager"`
	ClientSecret string `help:"Keycloak client secret" env:"KEYCLOAK_CLIENT_SECRET" required:""`
	Username     string `arg:"" name:"username" help:"Username for the new Keycloak user" required:""`
	Password     string `arg:"" name:"password" help:"Password for the new Keycloak user" required:""`
}

// AuthResponse represents Keycloak token response
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// Function to get Keycloak token using a client (not admin-cli)
func getClientToken(cli *CLI) (string, error) {
	resp, err := resty.R().
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     cli.ClientID,
			"client_secret": cli.ClientSecret,
		}).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Post(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", cli.KeycloakURL, cli.Realm))

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", resp.String())
	}

	var authResp AuthResponse
	if err := json.Unmarshal(resp.Body(), &authResp); err != nil {
		return "", err
	}

	return authResp.AccessToken, nil
}

// Function to create a new user in Keycloak
func createUser(cli *CLI, token string) (string, error) {
	resp, err := resty.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"username": cli.Username,
			"enabled":  true,
		}).
		Post(fmt.Sprintf("%s/admin/realms/%s/users", cli.KeycloakURL, cli.Realm))

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("failed to create user: %s", resp.String())
	}

	location := resp.Header().Get("Location")
	userID := location[strings.LastIndex(location, "/")+1:]

	return userID, nil
}

// Function to set user password (for updates, uses /reset-password)
func setUserPassword(cli *CLI, token, userID string) error {
	resp, err := resty.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"type":      "password",
			"value":     cli.Password,
			"temporary": false, // User keeps this password
		}).
		Put(fmt.Sprintf("%s/admin/realms/%s/users/%s/reset-password", cli.KeycloakURL, cli.Realm, userID))

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("failed to set password: %s", resp.String())
	}

	return nil
}

// Function to set an initial password using /credentials
func setInitialPassword(cli *CLI, token, userID string) error {
	resp, err := resty.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"type":      "password",
			"value":     cli.Password,
			"temporary": true, // User must change this password on first login
		}).
		Post(fmt.Sprintf("%s/admin/realms/%s/users/%s/credentials", cli.KeycloakURL, cli.Realm, userID))

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to set initial password: %s", resp.String())
	}

	return nil
}

func main() {
	// Parse CLI arguments and environment variables using Kong
	var cli CLI
	kong.Parse(&cli)

	// Get access token
	token, err := getClientToken(&cli)
	if err != nil {
		log.Fatalf("Error getting client token: %v", err)
	}

	// Create user
	userID, err := createUser(&cli, token)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}

	// Set initial temporary password
	err = setInitialPassword(&cli, token, userID)
	if err != nil {
		log.Fatalf("Error setting initial password: %v", err)
	}

	fmt.Println("User created successfully with an initial password (must be changed on first login)!")
}

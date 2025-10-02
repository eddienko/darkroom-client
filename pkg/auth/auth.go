package auth

import (
	"darkroom/pkg/config"
	"darkroom/pkg/netutil"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

func Login(cfg *config.Config, username string, password []byte) error {
	// TODO: Call real auth service
	// fmt.Printf("Authenticating %s against %s...\n", username, cfg.APIEndpoint)
	debug := config.Debug

	// Step 1: Login with username and password
	loginResp := sendLoginRequest(username, string(password), config.Debug)
	if *debug {
		fmt.Println("Login successful, userID:", loginResp.User.ID)
	}

	// Step 2: Prompt for OTP
	otp := promptPassword("Access Token: ")

	// Step 3: Validate OTP
	loginResponse := validateOTP(loginResp.User.ID, string(otp), debug)

	if loginResponse.AccessToken == "" {
		fmt.Println("Invalid access token. Please check your credentials and try again.")
		os.Exit(1)
	}

	cfg.AuthToken = loginResponse.AccessToken

	// Fetch additional credentials and user info
	kubeconfigContent, err := netutil.FetchKubeconfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to fetch kubeconfig: %w", err)
	}

	userInfo := GetUserInfo(cfg.AuthToken)

	cfg.KubeConfig = kubeconfigContent
	cfg.S3AccessToken = userInfo.S3AccessToken
	cfg.APIEndpoint = config.BaseURL
	cfg.UserName = userInfo.Username
	cfg.UserId = userInfo.ID

	if *debug {
		fmt.Println("Fetched kubeconfig and user info successfully.")
	}

	return nil
}

func sendLoginRequest(username, password string, debug *bool) LoginResponse {
	payload := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)

	resp := sendCustomPostRequest(config.LoginURL, payload, debug)
	defer resp.Body.Close()

	if *debug {
		fmt.Println("Login Response Status:", resp.Status)
		fmt.Println("Login Response Headers:", resp.Header)
	}

	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		exitWithError("Error decoding login response:", err)
	}

	if loginResponse.AccessToken == "" {
		fmt.Println("Error: Invalid username or password")
		os.Exit(1)
	}

	return loginResponse
}

func sendCustomPostRequest(url, payload string, debug *bool) *http.Response {
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		exitWithError("Error creating request:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		exitWithError("Error sending request:", err)
	}

	return resp
}

func validateOTP(uuid, token string, debug *bool) LoginResponse {
	payload := fmt.Sprintf(`{"uuid": "%s", "token": "%s"}`, uuid, token)

	resp := sendCustomPostRequest(config.ValidateOTPURL, payload, debug)
	defer resp.Body.Close()

	if *debug {
		fmt.Println("OTP Validation Response Status:", resp.Status)
		fmt.Println("OTP Validation Response Headers:", resp.Header)
	}

	// Parse response for any future validation or error handling
	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		exitWithError("Error decoding OTP validation response:", err)
	}

	return loginResponse

}

func GetUserInfo(token string) UserInfo {
	debug := config.Debug
	if token == "" {
		exitWithError("Error: No auth token found. Please login first.", nil)
	}
	req, err := http.NewRequest("GET", config.AboutMeURL, nil)
	if err != nil {
		exitWithError("Error creating request:", err)
	}

	// Set Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		exitWithError("Error sending request:", err)
	}
	defer resp.Body.Close()

	if *debug {
		fmt.Println("User Info Response Status:", resp.Status)
		fmt.Println("User Info Response Headers:", resp.Header)
	}

	// Decode response (adjust structure as needed)
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		exitWithError("Error decoding user info:", err)
	}

	// Print user info for debug
	if *debug {
		userJson, _ := json.MarshalIndent(userInfo, "", "  ")
		fmt.Println("User Info:", string(userJson))
	}

	return userInfo
}

func exitWithError(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
}

func promptPassword(prompt string) []byte {
	fmt.Print(prompt)
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		exitWithError("Error reading password:", err)
	}
	return pass
}

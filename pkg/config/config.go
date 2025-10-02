package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var debug bool = os.Getenv("DARKROOM_DEBUG") == "true"
var hardcodedKey = "&1cq^f_5ab7-$3yc-b(^$7t_=_c_@0gt+r2^%3mzee6jsaje-t"

var (
	BaseURL        = getBaseURL()
	LoginURL       = BaseURL + "/api/v1/accounts/token/"
	ValidateOTPURL = BaseURL + "/api/v1/accounts/otp/validate/"
	AboutMeURL     = BaseURL + "/api/v1/accounts/me/"
	StudyURL       = BaseURL + "/api/v1/archive/study/"
	ProjectURL     = BaseURL + "/api/v1/archive/project/"
	DatasetURL     = BaseURL + "/api/v1/archive/dataset/"
	KubeConfigURL  = BaseURL + "/api/v1/accounts/kubeconfig/"
	Debug          = &debug
)

func getBaseURL() string {
	if url := os.Getenv("DARKROOM_BASE_URL"); url != "" {
		return url
	}
	// Fallback default if env variable not set
	return "https://darkroom.ast.cam.ac.uk"
}

type Config struct {
	APIEndpoint   string `yaml:"apiEndpoint"`
	KubeConfig    string `yaml:"kubeConfig"`
	AuthToken     string `yaml:"authToken"`
	S3AccessToken string `yaml:"s3AccessToken"`
	UserName      string `yaml:"username"`
	UserId        int    `yaml:"userId"`
}

// New returns a default config
func New() *Config {
	return &Config{
		APIEndpoint: BaseURL,
	}
}

// getMachineID derives a machine-dependent ID (first non-loopback MAC).
func getMachineID() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		if len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("no suitable MAC address found")
}

// deriveFinalKey combines hardcoded key + machine ID into a 32-byte AES key.
func deriveFinalKey() ([]byte, error) {
	machinePart, err := getMachineID()
	if err != nil {
		return nil, err
	}
	combined := hardcodedKey + machinePart
	hash := sha256.Sum256([]byte(combined))
	return hash[:], nil
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".darkroom", "config.yaml.enc"), nil
}

// encrypt data with AES-GCM
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// decrypt data with AES-GCM
func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Load loads config from encrypted disk or returns default if not found
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// file not found → return default
		return New(), nil
	}

	key, err := deriveFinalKey()
	if err != nil {
		return nil, err
	}

	plaintext, err := decrypt(data, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(plaintext, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the config to encrypted disk
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	key, err := deriveFinalKey()
	if err != nil {
		return err
	}

	ciphertext, err := encrypt(data, key)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, ciphertext, 0600); err != nil {
		return err
	}

	if *Debug {
		fmt.Printf("Config saved to %s (encrypted)\n", path)
	}
	return nil
}

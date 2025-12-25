package services

import (
	"brute-force-login/config"
	"brute-force-login/database"
	"fmt"
	"os"
	"testing"
)

func TestUserSuspensionLogic(t *testing.T) {
	// Setup test database connection
	cfg := &config.Config{
		DBHost:            getEnv("TEST_DB_HOST", "localhost"),
		DBPort:            getEnv("TEST_DB_PORT", "5432"),
		DBUser:            getEnv("TEST_DB_USER", "postgres"),
		DBPassword:        getEnv("TEST_DB_PASSWORD", "postgres"),
		DBName:            getEnv("TEST_DB_NAME", "brute_force_login_test"),
		DBSSLMode:         "disable",
		UserAttemptLimit:  5,
		UserWindowMinutes: 5,
		UserSuspensionMinutes: 15,
		IPAttemptLimit:    100,
		IPWindowMinutes:   5,
	}

	err := database.InitDB(cfg)
	if err != nil {
		t.Skipf("Skipping test - database not available: %v", err)
		return
	}
	defer database.CloseDB()

	authService := NewAuthService(cfg)
	email := "test@example.com"
	ipAddress := "127.0.0.1"

	// Test: User should be suspended after 5 failed attempts
	for i := 0; i < 5; i++ {
		_, err := authService.Login(email, "wrongpassword", ipAddress)
		if err != nil {
			t.Logf("Attempt %d: %v", i+1, err)
		}
	}

	// 6th attempt should result in suspension
	response, err := authService.Login(email, "wrongpassword", ipAddress)
	if err != nil {
		t.Fatalf("Login should not return error: %v", err)
	}

	if response.Success {
		t.Error("Login should fail after 5 failed attempts")
	}

	if response.Message != "Account temporarily suspended due to too many failed attempts." {
		t.Errorf("Expected suspension message, got: %s", response.Message)
	}
}

func TestIPBlockLogic(t *testing.T) {
	// Setup test database connection
	cfg := &config.Config{
		DBHost:            getEnv("TEST_DB_HOST", "localhost"),
		DBPort:            getEnv("TEST_DB_PORT", "5432"),
		DBUser:            getEnv("TEST_DB_USER", "postgres"),
		DBPassword:        getEnv("TEST_DB_PASSWORD", "postgres"),
		DBName:            getEnv("TEST_DB_NAME", "brute_force_login_test"),
		DBSSLMode:         "disable",
		UserAttemptLimit:  5,
		UserWindowMinutes: 5,
		UserSuspensionMinutes: 15,
		IPAttemptLimit:    100,
		IPWindowMinutes:   5,
	}

	err := database.InitDB(cfg)
	if err != nil {
		t.Skipf("Skipping test - database not available: %v", err)
		return
	}
	defer database.CloseDB()

	authService := NewAuthService(cfg)
	ipAddress := "192.168.1.100"

	// Test: IP should be blocked after 100 failed attempts
	for i := 0; i < 100; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		_, err := authService.Login(email, "wrongpassword", ipAddress)
		if err != nil {
			t.Logf("Attempt %d: %v", i+1, err)
		}
	}

	// 101st attempt should result in IP block
	response, err := authService.Login("test@example.com", "wrongpassword", ipAddress)
	if err != nil {
		t.Fatalf("Login should not return error: %v", err)
	}

	if response.Success {
		t.Error("Login should fail after 100 failed attempts from same IP")
	}

	if response.Message != "IP temporarily blocked due to excessive failed login attempts." {
		t.Errorf("Expected IP block message, got: %s", response.Message)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


package services

import (
	"brute-force-login/config"
	"brute-force-login/database"
	"brute-force-login/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

func (s *AuthService) Login(email, password, ipAddress string) (*models.LoginResponse, error) {
	fmt.Println("Login attempt:", email, "from IP:", ipAddress, "password:", password)
	// Check IP block first
	isBlocked, err := s.isIPBlocked(ipAddress)
	if err != nil {
		return nil, fmt.Errorf("error checking IP block: %w", err)
	}
	if isBlocked {
		return &models.LoginResponse{
			Success: false,
			Message: "IP temporarily blocked due to excessive failed login attempts.",
		}, nil
	}

	// Check user suspension
	isSuspended, err := s.isUserSuspended(email)
	if err != nil {
		return nil, fmt.Errorf("error checking user suspension: %w", err)
	}
	if isSuspended {
		return &models.LoginResponse{
			Success: false,
			Message: "Account temporarily suspended due to too many failed attempts.",
		}, nil
	}

	// Get user from database
	user, err := s.getUserByEmail(email)
	fmt.Println("Retrieved user:", user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User doesn't exist - record failed attempt
			s.recordFailedAttempt(nil, email, ipAddress)
			return &models.LoginResponse{
				Success: false,
				Message: "Invalid email or password.",
			}, nil
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	fmt.Println("Comparing password for user:", user.Email)
	fmt.Println("Stored hash:", user.PasswordHash)
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	fmt.Println("err:", err)
	if err != nil {
		// Wrong password - record failed attempt
		s.recordFailedAttempt(&user.ID, email, ipAddress)
		return &models.LoginResponse{
			Success: false,
			Message: "Invalid email or password.",
		}, nil
	}

	// Successful login - generate token
	token, err := s.generateToken(user.ID, email)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &models.LoginResponse{
		Success: true,
		Message: "Login successful.",
		Token:   token,
	}, nil
}

func (s *AuthService) getUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	row := database.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) recordFailedAttempt(userID *int, email, ipAddress string) error {
	// Record user-level failed attempt
	query := `INSERT INTO user_failed_attempts (user_id, email, ip_address) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, userID, email, ipAddress)
	if err != nil {
		return fmt.Errorf("error recording user failed attempt: %w", err)
	}

	// Record IP-level failed attempt
	query = `INSERT INTO ip_failed_attempts (ip_address, email) VALUES ($1, $2)`
	_, err = database.DB.Exec(query, ipAddress, email)
	if err != nil {
		return fmt.Errorf("error recording IP failed attempt: %w", err)
	}

	// Check if user should be suspended
	if userID != nil {
		count, err := s.getUserFailedAttemptCount(email)
		if err != nil {
			return err
		}
		if count >= s.cfg.UserAttemptLimit {
			err = s.suspendUser(*userID, email)
			if err != nil {
				return err
			}
		}
	}

	// Check if IP should be blocked
	count, err := s.getIPFailedAttemptCount(ipAddress)
	fmt.Println("IP failed attempt count for", ipAddress, "is", count)
	if err != nil {
		return err
	}
	if count >= s.cfg.IPAttemptLimit {
		err = s.blockIP(ipAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AuthService) getUserFailedAttemptCount(email string) (int, error) {
	windowStart := time.Now().Add(-time.Duration(s.cfg.UserWindowMinutes) * time.Minute)
	query := `SELECT COUNT(*) FROM user_failed_attempts 
		WHERE email = $1 AND attempted_at > $2`

	var count int
	err := database.DB.QueryRow(query, email, windowStart).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting user failed attempts: %w", err)
	}
	return count, nil
}

func (s *AuthService) getIPFailedAttemptCount(ipAddress string) (int, error) {
	windowStart := time.Now().Add(-time.Duration(s.cfg.IPWindowMinutes) * time.Minute)
	query := `SELECT COUNT(*) FROM ip_failed_attempts 
		WHERE ip_address = $1 AND attempted_at > $2`

	var count int
	err := database.DB.QueryRow(query, ipAddress, windowStart).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting IP failed attempts: %w", err)
	}
	return count, nil
}

func (s *AuthService) isUserSuspended(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM user_suspensions 
		WHERE email = $1 AND suspended_until > NOW()`

	var count int
	err := database.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking user suspension: %w", err)
	}
	return count > 0, nil
}

func (s *AuthService) suspendUser(userID int, email string) error {
	suspendedUntil := time.Now().Add(time.Duration(s.cfg.UserSuspensionMinutes) * time.Minute)
	query := `INSERT INTO user_suspensions (user_id, email, suspended_until) 
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`

	_, err := database.DB.Exec(query, userID, email, suspendedUntil)
	if err != nil {
		return fmt.Errorf("error suspending user: %w", err)
	}
	return nil
}

func (s *AuthService) isIPBlocked(ipAddress string) (bool, error) {
	query := `SELECT COUNT(*) FROM ip_blocks 
		WHERE ip_address = $1 AND blocked_until > NOW()`

	var count int
	err := database.DB.QueryRow(query, ipAddress).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking IP block: %w", err)
	}
	return count > 0, nil
}

func (s *AuthService) blockIP(ipAddress string) error {
	blockedUntil := time.Now().Add(time.Duration(s.cfg.IPWindowMinutes) * time.Minute)
	query := `INSERT INTO ip_blocks (ip_address, blocked_until) 
		VALUES ($1, $2)
		ON CONFLICT (ip_address) 
		DO UPDATE SET blocked_until = EXCLUDED.blocked_until`

	_, err := database.DB.Exec(query, ipAddress, blockedUntil)
	if err != nil {
		return fmt.Errorf("error blocking IP: %w", err)
	}
	return nil
}

func (s *AuthService) generateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

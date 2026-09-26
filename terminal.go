package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Terminal struct {
	isAuthorized bool
	passwordHash []byte
}

func NewTerminal() *Terminal {
	return &Terminal{
		isAuthorized: false,
	}
}

func (t *Terminal) LoadPasswordHash(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		defaultHash := []byte("$2a$10$R9h/azU1Z7E.pX.0gQ8oPefA4sB4r2W3k3o6lCqQx6t.W8pQ4z82e")
		_ = os.WriteFile(filePath, defaultHash, 0644)
		t.passwordHash = defaultHash
		return fmt.Errorf("configuration file missing; generated default fallback file instead")
	}
	t.passwordHash = bytes.TrimSpace(data)
	return nil
}

func (t *Terminal) Execute(rawInput string) string {
	cleanedInput := strings.TrimSpace(rawInput)
	normalizedCommand := strings.ToLower(cleanedInput)

	if normalizedCommand == "exit" {
		go func() { os.Exit(0) }()
		return "SHUTDOWN"
	}

	if !t.isAuthorized {
		return t.handleAuthentication(cleanedInput)
	}

	return t.processSystemCommand(normalizedCommand, cleanedInput)
}

func (t *Terminal) handleAuthentication(passwordAttempt string) string {
	err := bcrypt.CompareHashAndPassword(t.passwordHash, []byte(passwordAttempt))
	if err != nil {
		return "AUTH_FAIL"
	}
	t.isAuthorized = true
	return "AUTH_SUCCESS"
}

func (t *Terminal) processSystemCommand(cmd string, original string) string {
	switch cmd {
	case "help":
		return "Available Commands:\n  help    - Display this menu\n  clear   - Clear terminal screen\n  status  - Query host system integrity\nlogout  - Revoke authentication level\n  exit    - Terminate server runtime and close card"
	case "status":
		return "[OK] Golang server cluster state nominal.\n[AUTH] Sub-level authorization verified: root."
	case "logout":
		t.isAuthorized = false
		return "LOGOUT"
	case "clear":
		return ""
	default:
		return fmt.Sprintf("Error: Command '%s' not recognized. Type 'help' for layout instructions.", original)
	}
}

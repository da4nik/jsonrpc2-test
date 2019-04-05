package auth

import "fmt"

const secretPass = "123456789"

// Authenticate - authenticates user
func Authenticate(email, password string) (string, error) {
	if password == secretPass {
		return "supermegatoken", nil
	}
	return "", fmt.Errorf("wrong email/password combination")
}

package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := map[string]string{
		"admin123":    "",
		"manager123":  "",
		"employee123": "",
	}

	for password := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error hashing password %s: %v", password, err)
		}
		passwords[password] = string(hash)
		fmt.Printf("Password: %s\nHash: %s\n\n", password, string(hash))
	}
}

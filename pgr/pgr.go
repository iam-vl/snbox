package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	pwd = "Mp3652#vl1876"
)

func main() {
	GenerateHash(pwd)

}

func GenerateHash(password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	fmt.Printf("Hash value: %+v\t\nHash type: %T\n", hash, hash)
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		fmt.Println("Passwords match")
	}
}

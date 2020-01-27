package main

import (
	"flag"

	"tourtoster/handler"
)

func main() {
	var password string
	flag.StringVar(&password, "password", "", "string to bcrypt hash")
	flag.Parse()

	if password == "" {
		return
	}

	passwordHash, err := handler.HashPassword(password)
	if err != nil {
		println("error create hash", "error", err.Error())
		return
	}

	println("Result", passwordHash)
}

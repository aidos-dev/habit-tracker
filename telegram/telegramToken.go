package app

import (
	"flag"
	"log"
)

/*
mustToken func name has prefix must indicating that the returned
value from this function is mandatory for the program.
If mustToken can't find a token it stops the process.
Prefix must allows us not to handle an error in case of token is missing.
It is done only to reduce the code amount
*/
func MustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token to access the telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("error: token is not specified")
	}

	return *token
}

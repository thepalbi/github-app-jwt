package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/urfave/cli/v2"
)

func main() {
	var privateKeyFile string
	var githubAppId int
	var debug bool

	app := cli.NewApp()
	app.Name = "github-app-jwt"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "private-key-file",
			Usage:       "Path to the private key file",
			Destination: &privateKeyFile,
			Required:    true,
		},
		&cli.IntFlag{
			Name:        "github-app-id",
			Usage:       "Github App ID",
			Destination: &githubAppId,
			Required:    true,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug logging",
			Destination: &debug,
		},
	}

	dlog := func(format string, args ...interface{}) {
		if debug {
			log.Printf(format, args...)
		}
	}

	app.Action = func(ctx *cli.Context) error {
		// Github App JWT according to
		// https: //docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app
		now := time.Now().Add(-time.Minute * 1) // fight again clock drift to be able to use token now
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"iat": now.Unix(),
			"exp": now.Add(time.Minute * 10).Unix(),
			"iss": githubAppId,
		})

		bs, err := os.ReadFile(privateKeyFile)
		if err != nil {
			dlog("Error reading PEM: %s", err)
			os.Exit(1)
		}

		pk, err := jwt.ParseRSAPrivateKeyFromPEM(bs)
		if err != nil {
			dlog("Error parsing PEM: %s", err)
			os.Exit(1)
		}

		signedToken, err := token.SignedString(pk)
		if err != nil {
			dlog("Error signing token: %s", err)
			os.Exit(1)
		}

		fmt.Print(signedToken)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

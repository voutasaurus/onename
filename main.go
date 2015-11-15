package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
)

var BASE_URL = "https://api.onename.com/v1"

type client struct {
	ApiID     string
	ApiSecret string
	BaseUrl   string
}

func getClient(useLive bool) (client, error) {
	file, err := os.Open("cred.txt")
	if err != nil {
		return client{}, err
	}

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return client{}, err
	}

	if useLive {
		live := "https://api.onename.com/v1"
		return client{lines[0], lines[1], live}, nil
	}
	test := "http://localhost:12345"
	return client{lines[0], lines[1], test}, nil
}

func oneNameClientApp(out io.Writer) *cli.App {
	app := cli.NewApp()
	app.Name = "onename"
	app.Usage = "CLI for the OneName API"
	app.Writer = out
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "live, l",
			Usage: "point to the live API",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "lookup",
			Aliases: []string{"ls"},
			Usage:   "lookup user records - onename lookup fredwilson albertwenger",
			Action: func(c *cli.Context) {
				httpClient, err := getClient(c.GlobalBool("live"))
				if err != nil {
					fmt.Fprintln(out, err)
					return
				}
				users, err := httpClient.GetUsers(c.Args())
				if err != nil {
					fmt.Fprintln(out, err)
					return
				}
				fmt.Fprintln(out, users)
			},
		},
		{
			Name:    "stats",
			Aliases: []string{"st"},
			Usage:   "get onename user stats - onename stats",
			Action: func(c *cli.Context) {
				httpClient, err := getClient(c.GlobalBool("live"))
				if err != nil {
					fmt.Fprintln(out, err)
					return
				}
				stats, err := httpClient.GetUserStats()
				if err != nil {
					fmt.Fprintln(out, err)
				} else {
					fmt.Fprintln(out, stats)
				}
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s", "find"},
			Usage:   "search for a user - onename search wenger",
			Action: func(c *cli.Context) {
				httpClient, err := getClient(c.GlobalBool("live"))
				if err != nil {
					fmt.Fprintln(out, err)
					return
				}
				searchResults, err := httpClient.SearchUsers(c.Args().First())
				if err != nil {
					fmt.Fprintln(out, err)
				} else {
					fmt.Fprintln(out, searchResults)
				}
			},
		},
		{
			Name:    "unspents",
			Aliases: []string{"un"},
			Usage:   "unspents of a bitcoin address - onename unspents 1QHDGGLEKK7FZWsBEL78acV9edGCTarqXt",
			Action: func(c *cli.Context) {
				httpClient, err := getClient(c.GlobalBool("live"))
				if err != nil {
					fmt.Fprintln(out, err)
					return
				}
				unspentsResults, err := httpClient.GetUnspents(c.Args().First())
				if err != nil {
					fmt.Fprintln(out, err)
				} else {
					fmt.Fprintln(out, unspentsResults)
				}
			},
		},
		{
			Name:    "names",
			Aliases: []string{"n"},
			Usage:   "names associated with a bitcoin address - onename names 1QHDGGLEKK7FZWsBEL78acV9edGCTarqXt",
			Action: func(c *cli.Context) {
				httpClient, err := getClient(c.GlobalBool("live"))
				if err != nil {
					fmt.Fprintln(out, err)
					return
				}
				namesResults, err := httpClient.GetNames(c.Args().First())
				if err != nil {
					fmt.Fprintln(out, err)
				} else {
					fmt.Fprintln(out, namesResults)
				}
			},
		},
		{
			Name:    "dkim",
			Aliases: []string{"d"},
			Usage:   "dkim info on a domain - onename dkim onename.com",
			Action: func(c *cli.Context) {
				httpClient, err := getClient(c.GlobalBool("live"))
				if err != nil {
					fmt.Fprintln(out, err)
					return
				}
				dkimResult, err := httpClient.GetDKIMInfo(c.Args().First())
				if err != nil {
					fmt.Fprintln(out, err)
				} else {
					fmt.Fprintln(out, dkimResult)
				}
			},
		},
	}
	app.Action = func(c *cli.Context) {

		return
	}

	return app

}

func main() {
	app := oneNameClientApp(os.Stdout)
	app.Run(os.Args)
}

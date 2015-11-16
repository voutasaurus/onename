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
				if len(c.Args()) > 0 {
					users, err := httpClient.GetUsers(c.Args())
					if err != nil {
						fmt.Fprintln(out, "Couldn't get users:", fmt.Sprint(c.Args()), "Got error:\n\t", err.Error())
						return
					}
					fmt.Fprintln(out, users)
				} else {
					allUsers, err := httpClient.GetAllUsers()
					if err != nil {
						fmt.Fprintln(out, "No user specified so attempted to get all users. Got error:\n\t", err.Error())
						return
					}
					fmt.Fprintln(out, allUsers)
				}
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
		{
			Name:    "register",
			Aliases: []string{"r"},
			Usage:   "register a name/address association - onename register name address",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 2 {
					fmt.Fprintln(out, "You must specify a name and an address to register.")
					return
				}
				httpClient, err := getClient(c.GlobalBool("live"))
				if err != nil {
					fmt.Fprintln(out, "Error constructing http client:", err.Error())
					return
				}
				ack, err := httpClient.RegisterUser(c.Args()[0], c.Args()[1])
				if err != nil {
					fmt.Fprintln(out, "Error registering user:\n\t", err)
				} else {
					fmt.Fprintln(out, ack)
				}
			},
		},
		{
			Name:    "broadcast",
			Aliases: []string{"transaction"},
			Usage:   "Broadcast a transaction - onename broadcast trasactionhex",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					fmt.Fprintln(out, "You must specify a transaction to broadcast.")
					return
				}
				httpClient, err := getClient(c.GlobalBool("live"))
				if err != nil {
					fmt.Fprintln(out, "Error constructing http client:", err.Error())
					return
				}
				ack, err := httpClient.BroadcastTransactions(c.Args().First())
				if err != nil {
					fmt.Fprintln(out, "Error broadcasting transaction:\n\t", err)
				} else {
					fmt.Fprintln(out, ack)
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

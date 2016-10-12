package main

import (
	"github.com/codegangsta/cli"
	"fmt"
	"os"
	"github.com/orieken/cheeseweb/fetcher"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "foo",
			Value: "Foo",
			Usage: "foo thing",
		},
	}

	app.Name = "cheeseweb"
	app.Usage = "Spin up a Selenium Grid"

	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		urls := []string{
			"https://selenium-release.storage.googleapis.com/2.53/selenium-server-standalone-2.53.1.jar",
			"https://chromedriver.storage.googleapis.com/2.24/chromedriver_mac64.zip",
			"https://github.com/mozilla/geckodriver/releases/download/v0.11.1/geckodriver-v0.11.1-macos.tar.gz",
		}
		fetcher.Fetch(urls)
		return nil
	}
	app.Run(os.Args)
}
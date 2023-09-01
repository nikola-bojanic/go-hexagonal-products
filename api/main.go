// Starts the CLI application, executes any given command, and shuts down
// Added to avoid having to install the egw code on the machine to execute commands
package main

import (
	"os"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/api/cmd"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/app"
	"github.com/urfave/cli"
)

func RunCLIApplication() {

	app := app.MustInitializeApp()

	cliApp := cli.NewApp()
	cliApp.Name = "exchange-gateway"
	cliApp.Description = "Command line utility for egw development"
	cliApp.Commands = []cli.Command{
		cmd.NewDbCmd(app),
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		app.Logger.Error("failed running command", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func main() {
	RunCLIApplication()
}

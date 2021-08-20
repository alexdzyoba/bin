package main

import (
	"log"
	"os"
	"path"

	"github.com/urfave/cli/v2"
)

var (
	version = "0.0.1"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("failed to get home directory:", err)
	}

	var config Config
	app := &cli.App{
		Name:      "bin",
		Usage:     "Binary programs manager",
		UsageText: "bin <command> [options] [arguments...]",
		Commands: []*cli.Command{
			{
				Name:      "install",
				Aliases:   []string{"i"},
				Usage:     "install binary",
				UsageText: "bin install [options] <URL>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "target.filename",
						Aliases:     []string{"o"},
						Usage:       "Override target filename",
						DefaultText: "filename from archive as is",
						EnvVars:     []string{"BIN_TARGET_FILENAME"},
						Destination: &config.Target.Filename,
					},
					&cli.StringFlag{
						Name:        "target.dir",
						Aliases:     []string{"d"},
						Usage:       "Override target directory",
						Value:       path.Join(home, "bin"),
						EnvVars:     []string{"BIN_TARGET_DIR"},
						Destination: &config.Target.Dir,
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						log.Println("Single URL argument is required")
						return cli.ShowCommandHelp(c, c.Command.Name)
					}

					return install(c.Args().Get(0), config)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list installed binaries",
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

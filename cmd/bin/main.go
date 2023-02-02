package main

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	version = "0.0.1"
)

func main() {
	verbosity := Verbosity{zerolog.InfoLevel}
	verbosity.Init()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get home directory")
	}

	var config Config
	app := &cli.App{
		Name:      "bin",
		Usage:     "Binary programs manager",
		UsageText: "bin <command> [options] [arguments...]",

		Flags: []cli.Flag{
			&cli.GenericFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Set verbosity level",
				EnvVars: []string{"BIN_VERBOSITY"},
				Value:   &verbosity,
			},
		},

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
						log.Error().Msg("single URL argument is required")
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
		log.Fatal().Err(err)
	}
}

type Verbosity struct {
	level zerolog.Level
}

func (v *Verbosity) Init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"}).With().Caller().Logger()
	zerolog.SetGlobalLevel(v.level)
}

func (v *Verbosity) Set(value string) error {
	level, err := zerolog.ParseLevel(value)
	if err != nil {
		return fmt.Errorf("unknown log level '%s'", value)
	}

	v.level = level
	zerolog.SetGlobalLevel(v.level)
	return nil
}

func (v *Verbosity) String() string {
	return v.level.String()
}

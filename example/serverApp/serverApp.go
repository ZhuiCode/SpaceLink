package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	ServerAddress string `long:"server-address" description:"server address, required"`
}

func mian() {
	var opt Options
	parser := flags.NewParser(&opt, flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	if opt.ServerAddress == "" {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

}

package main

import (
	"flag"
	"os"

	"github.com/spf13/viper"
)

// Flags are the flags of the program.
type Flags struct {
	ListenAddress string
	Debug         bool
	CertFile      string
	KeyFile       string
}

// NewFlags returns the flags of the commandline
func NewFlags(cfg *viper.Viper) *Flags {
	flags := &Flags{}
	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&flags.ListenAddress, "listen-address", cfg.Get("serverPort").(string), "webhook server listen address")
	fl.BoolVar(&flags.Debug, "debug", cfg.Get("mode").(bool), "enable debug mode")
	fl.StringVar(&flags.CertFile, "tls-cert-file", cfg.Get("certFile").(string), "TLS certificate file")
	fl.StringVar(&flags.KeyFile, "tls-key-file", cfg.Get("keyFile").(string), "TLS key file")
	fl.Parse(os.Args[1:])

	return flags
}

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zqzca/back/app"
)

var cdn string
var secure bool
var livereload bool
var bindhttp string
var bindscp string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "zqz",
		Short: "the entire backend for zqz.ca.",
		Long:  ``,
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Runs ZQZ",
		Long:  "Runs ZQZ",

		Run: func(cmd *cobra.Command, args []string) {
			cfg := app.Config{
				Secure:       secure,
				LiveReload:   livereload,
				CDNURL:       cdn,
				HTTPBindAddr: bindhttp,
				SCPBindAddr:  bindscp,
			}

			app.Run(cfg)
		},
	}

	rootCmd.AddCommand(serveCmd)

	serveFlags := serveCmd.Flags()
	serveFlags.BoolVar(&secure, "secure", false, "Serve HTTP2 instead of HTTP")
	serveFlags.BoolVar(&livereload, "livereload", false, "Enable LiveReload")
	serveFlags.StringVar(&cdn, "cdn", "/assets", "URL for assets")
	serveFlags.StringVar(&bindhttp, "http", ":3001", "HTTP Bind address")
	serveFlags.StringVar(&bindscp, "scp", ":2020", "SCP Bind address")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

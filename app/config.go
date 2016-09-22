package app

// Config contains all settings required to start zqz.
type Config struct {
	HTTPBindAddr string
	SCPBindAddr  string
	LiveReload   bool
	CDNURL       string
	Secure       bool
}

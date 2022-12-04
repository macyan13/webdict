package server

import "time"

// Opts flags and envs to run server
type Opts struct {
	Auth  AuthGroup  `group:"auth" namespace:"auth" env-namespace:"AUTH"`
	Admin AdminGroup `group:"admin" namespace:"admin" env-namespace:"ADMIN"`

	Port       int    `long:"port" env:"PORT" default:"4000" description:"port"`
	WebdictURL string `long:"url" env:"URL" description:"url to webdict"`
	Dbg        bool   `long:"dbg" env:"DEBUG" description:"debug mode"`
}

// AuthGroup defines options group for auth params
type AuthGroup struct {
	TTL struct {
		Auth    time.Duration `long:"auth" env:"AUTH" default:"10m" description:"auth JWT TTL"`
		Refresh time.Duration `long:"refresh" env:"REFRESH" default:"24h" description:"refresh JWT TTL"`
		Cookie  time.Duration `long:"cookie" env:"COOKIE" default:"200h" description:"refresh cookie TTL"`
	} `group:"ttl" namespace:"ttl" env-namespace:"TTL"`

	Secret string `long:"secret" env:"SECRET" required:"true" description:"the secret key used to sign JWT, should be a random, long, hard-to-guess string"`
}

// AdminGroup defines options group for admin user params
type AdminGroup struct {
	AdminPasswd string `long:"passwd" env:"PASSWD" default:"admin_password" description:"admin user password"`
	AdminEmail  string `long:"email" env:"EMAIL" default:"admin@email.com" description:"admin user email"`
}

package server

import "time"

// Opts flags and envs to run server
type Opts struct {
	Auth  AuthGroup  `group:"auth" namespace:"auth" env-namespace:"AUTH"`
	Admin AdminGroup `group:"admin" namespace:"admin" env-namespace:"ADMIN"`
	Mongo MongoGroup `group:"mongo" namespace:"mongo" env-namespace:"MONGO"`
	Cache CacheGroup `group:"cache" namespace:"cache" env-namespace:"CACHE"`

	Port       int    `long:"port" env:"PORT" default:"4000" description:"port"`
	WebdictURL string `long:"url" env:"URL" description:"url to webdict"`
	Dbg        bool   `long:"dbg" env:"DEBUG" description:"debug mode"`
}

// AuthGroup defines options group for auth params
type AuthGroup struct {
	TTL struct {
		Auth    time.Duration `long:"auth" env:"AUTH" default:"2h" description:"auth JWT TTL"`
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

// MongoGroup defines options group for mongo connection
type MongoGroup struct {
	Database string `long:"database" env:"DB" default:"webdict" description:"name of the mongo DB"`
	Host     string `long:"host" env:"HOST" default:"localhost" description:"mongo DB host"`
	Port     int    `long:"port" env:"PORT" default:"27017" description:"mongo DB port"`
	Username string `long:"username" env:"USERNAME" default:"root" description:"name of the mongo username"`
	Passwd   string `long:"password" env:"PASSWD" default:"example" description:"mongo password"`
}

// CacheGroup defines options group for in memory cache
type CacheGroup struct {
	TagCacheTTL         time.Duration `long:"tag_cache_ttl" env:"TAG_CACHE_TTL" default:"3600s" description:"Cache TTL for tags"`
	TranslationCacheTTL time.Duration `long:"translation_cache_ttl" env:"TRANSLATION_CACHE_TTL" default:"3600s" description:"Cache TTL for translations"`
	LangCacheTTL        time.Duration `long:"lang_cache_ttl" env:"LANG_CACHE_TTL" default:"3600s" description:"Cache TTL for languages"`
}

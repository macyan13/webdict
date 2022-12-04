package server

import "time"

func initTestServer() *HttpServer {
	authGroup := AuthGroup{}
	authGroup.TTL.Auth = time.Minute * 10
	authGroup.TTL.Refresh = time.Minute * 10
	authGroup.TTL.Cookie = time.Hour

	return InitServer(Opts{
		Auth: authGroup,
		Admin: AdminGroup{
			AdminPasswd: "test_password",
			AdminEmail:  "test@email.com",
		},
		Port:       4000,
		WebdictURL: "",
		Dbg:        false,
	})
}

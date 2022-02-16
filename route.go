package main

import "github.com/PengLei-Adam/lade/framework"

func registerRouter(core *framework.Core) {
	core.Get("/foo", FooControllerHandler)
	core.Get("/login", UserLoginController)
}

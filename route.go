package main

import "coredemo/framework"

func registerRouter(core *framework.Core) {
	core.Get("/foo", FooControllerHandler)
	core.Get("/login", UserLoginController)
}

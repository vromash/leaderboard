package main

import (
	grpc_controller "main/internal/controller/grpc"
	"main/internal/state/config"
	"main/internal/state/ctx"
)

func main() {
	config.SetupConfig()
	config.SetupLogger()

	appCtx := ctx.Setup()

	grpc_controller.RunServer(appCtx)
}

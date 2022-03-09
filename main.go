package leaderboard

import (
	"main/internal/state/config"
	"main/internal/state/ctx"
)

func main() {
	config.SetupConfig()
	config.SetupLogger()

	appCtx := ctx.Setup()
}

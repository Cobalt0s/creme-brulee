package harvester

import (
	"context"
	"fmt"
	"github.com/Cobalt0s/creme-brulee/pkg/rest/gintonic"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
)

func Prompt(ctx context.Context) {
	envHarvesterHost, err := config.GetEnv("HARVESTER_HOST")
	if err != nil {
		log := ctxlogrus.Extract(ctx)
		log.Warn("Harvester prompt failed, missing HARVESTER_HOST env var")
		return
	}
	_, _ = gintonic.MakeRequest(ctx, "POST", fmt.Sprintf("%v/prompt", envHarvesterHost), nil)
}

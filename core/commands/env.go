package commands

import (
	"fmt"

	"github.com/AtlantPlatform/go-ipfs/commands"
	"github.com/AtlantPlatform/go-ipfs/core"
	"github.com/AtlantPlatform/go-ipfs/repo/config"
)

// GetNode extracts the node from the environment.
func GetNode(env interface{}) (*core.IpfsNode, error) {
	ctx, ok := env.(*commands.Context)
	if !ok {
		return nil, fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
	}

	return ctx.GetNode()
}

// GetConfig extracts the config from the environment.
func GetConfig(env interface{}) (*config.Config, error) {
	ctx, ok := env.(*commands.Context)
	if !ok {
		return nil, fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
	}

	return ctx.GetConfig()
}

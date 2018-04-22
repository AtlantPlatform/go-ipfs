package commands

import (
	cmds "github.com/AtlantPlatform/go-ipfs/commands"
	"github.com/AtlantPlatform/go-ipfs/go-ipfs-cmdkit"
)

var DiagCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Generate diagnostic reports.",
	},

	Subcommands: map[string]*cmds.Command{
		"sys":  sysDiagCmd,
		"cmds": ActiveReqsCmd,
	},
}

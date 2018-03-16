package commands

import (
	cmds "bitbucket.org/atlantproject/go-ipfs/commands"
	"bitbucket.org/atlantproject/go-ipfs/go-ipfs-cmdkit"
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

package commands

import (
	"errors"

	cmds "bitbucket.org/atlantproject/go-ipfs/commands"
	cmdkit "bitbucket.org/atlantproject/go-ipfs/go-ipfs-cmdkit"
)

var MountCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "Not yet implemented on Windows.",
		ShortDescription: "Not yet implemented on Windows. :(",
	},

	Run: func(req cmds.Request, res cmds.Response) {
		res.SetError(errors.New("Mount isn't compatible with Windows yet"), cmdkit.ErrNormal)
	},
}

package command

import (
	"github.com/baiyecha/cloud_disk/server"
	"github.com/urfave/cli"
)

func RegisterCommand(svr *server.Server) []cli.Command {
	return []cli.Command{
		NewExampleCommand(svr),
	}
}

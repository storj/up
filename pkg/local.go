package sjr

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
)

func localCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "local",
		Short: "Bind mount local executable to run it inside the container",
		RunE: func(cmd *cobra.Command, args []string) error {
			return UpdateEach(service, Local)
		},
	})
}

func Local(service *ServiceConfig) error {
	cmd := service.Command[0]
	goBinPath := path.Join(os.Getenv("GOPATH"), "bin")
	service.Volumes = append(service.Volumes, fmt.Sprintf(
		"%s:%s",
		path.Join(goBinPath, cmd),
		path.Join("/var/lib/storj/go/bin", cmd)))
	return nil
}

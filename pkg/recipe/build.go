package recipe

import (
	_ "embed"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

//go:embed edge.Dockerfile
var edgeDocker []byte

//go:embed storj.Dockerfile
var storjDocker []byte

func buildCmd(service string, command *cobra.Command) {
	var repository, branch string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build image on-the-fly instead of using pre-baked image",
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.UpdateEach(service, func(service *common.ServiceConfig) error {
				return Build(service, repository, branch)
			})
		},
	}

	cmd.Flags().StringVarP(&repository, "repository", "r", "https://github.com/storj/{gateway-mt/storj}.git", "Git repository to clone before build.")
	cmd.Flags().StringVarP(&branch, "branch", "b", "main", "The branch to checkout and build")
	command.AddCommand(cmd)
}

func Build(service *common.ServiceConfig, repository string, branch string) error {
	imageType := "storj"
	err := common.ExtractFile(edgeDocker, "edge.Dockerfile")
	if err != nil {
		return err
	}

	err = common.ExtractFile(storjDocker, "storj.Dockerfile")
	if err != nil {
		return err
	}
	if strings.Contains(service.Image, "-edge") {
		imageType = "edge"
		repository = strings.ReplaceAll(repository, "{gateway-mt/storj}", "gateway-mt")
	} else {
		repository = strings.ReplaceAll(repository, "{gateway-mt/storj}", "storj")
	}
	service.Image = ""
	service.Build = &common.BuildConfig{
		Context:    ".",
		Dockerfile: imageType + ".Dockerfile",
		Args: map[string]*string{
			"REPO":   &repository,
			"BRANCH": &branch,
		},
	}
	return nil
}

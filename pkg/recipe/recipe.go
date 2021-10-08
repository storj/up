package recipe

import "github.com/spf13/cobra"

func AddMutatorCommands(group string, serviceCmd *cobra.Command) {
	debugCmd(group, serviceCmd)
	imageCmd(group, serviceCmd)
	localCmd(group, serviceCmd)
	versionCmd(group, serviceCmd)
	localEntrypointCmd(group, serviceCmd)
	envCmd(group, serviceCmd)
	buildCmd(group, serviceCmd)
	scaleCmd(group, serviceCmd)
}


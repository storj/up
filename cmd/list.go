package cmd

import (
	"fmt"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all the configured services from the dockerfile",
	RunE: func(cmd *cobra.Command, args []string) error {
		return List()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func List() error {
	project, err := common.CreateComposeProject("docker-compose.yaml")
	if err != nil {
		return err
	}

	for _, service := range project.AllServices() {
		fmt.Println(service.Name, service.Image)
	}
	return nil
}
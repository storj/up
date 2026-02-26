// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/compose"
	"storj.io/storj-up/pkg/runtime/runtime"
)

var cfgFile string
var rootDir string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "storj-up",
	Args:  cobra.NoArgs,
	Short: "A golang wrapper for creating customized Storj environment with many components",
	Long: `storj-up can be used to create containerized or standalone runtime environment,
leverages existing binaries, or even references code to be built in docker and create images.

For example:

storj-up build remote gerrit 5826`,
}

// Execute is the execution of the top level storj-up command.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatalf("%++v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.reminderctl.yaml)")
	RootCmd.PersistentFlags().StringVar(&rootDir, "root", "", "The directory of the project. If not set, the current directory is used.")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".sjr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sjr")
	}

	viper.SetEnvPrefix("STORJUP")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// ExecuteStorjUP can execute any operation with loaded stack/runtime and write back the results.
func ExecuteStorjUP(exec func(stack recipe.Stack, rt runtime.Runtime, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if rootDir != "" {
			pwd = rootDir
		}
		rt, err := FromDir(pwd)
		if err != nil {
			return err
		}
		st, err := recipe.GetStack()
		if err != nil {
			return err
		}
		err = rt.Reload(st)
		if err != nil {
			return err
		}

		err = exec(st, rt, args)
		if err != nil {
			return err
		}
		return rt.Write()
	}
}

// ChangeCompose applies modification to compose based runtime services. Used mainly in legacy commands.
func ChangeCompose(st recipe.Stack, rt runtime.Runtime, selectors []string, do func(composeService *types.ServiceConfig) error) error {
	return runtime.ModifyService(st, rt, selectors, func(s runtime.Service) error {
		composeService, ok := s.(*compose.Service)
		if !ok {
			return errs.Errorf("this subcommand is supported only for compose based environments")
		}
		return composeService.TransformRaw(do)
	})
}

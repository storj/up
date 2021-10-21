package sjr

import (
	"fmt"
	"github.com/spf13/cobra"
	"reflect"
	"regexp"
	"storj.io/storj/storagenode"
	"strings"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use: "configs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printConfigs()
		},
	})
}

func printConfigs() error {
	configType := reflect.TypeOf(storagenode.Config{})
	printConfigStruct("STORJ", configType)
	return nil
}

func printConfigStruct(prefix string, configType reflect.Type) {
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			printConfigStruct(prefix+"_"+upper(field.Name), field.Type)
		} else {
			fmt.Println(prefix + "_" + upper(field.Name) + " " + field.Tag.Get("help") + " " + field.Tag.Get("default"))
		}
	}

}

func upper(name string) string {
	smallCapital := regexp.MustCompile("([a-z])([A-Z])")
	name = smallCapital.ReplaceAllString(name, "${1}_$2")
	return strings.ToUpper(name)
}

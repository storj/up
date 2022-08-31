// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/zeebo/errs"

	"storj.io/gateway-mt/pkg/auth"
	"storj.io/gateway-mt/pkg/linksharing"
	"storj.io/storj/satellite"
	"storj.io/storj/storagenode"
)

var configTypes = map[string]reflect.Type{
	"storagenode":     reflect.TypeOf(storagenode.Config{}),
	"satellite-api":   reflect.TypeOf(satellite.Config{}),
	"satellite-admin": reflect.TypeOf(satellite.Config{}),
	"satellite-core":  reflect.TypeOf(satellite.Config{}),
	"linksharing":     reflect.TypeOf(linksharing.Config{}),
	"authservice":     reflect.TypeOf(auth.Config{}),
}

func main() {
	for name, t := range configTypes {
		err := generate(name, t)
		if err != nil {
			panic(err)
		}
	}
}

func generate(name string, t reflect.Type) error {
	fileName := name + ".go"
	fmt.Println("Writing " + fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = f.WriteString(`// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.
package config

func init() {
`)
	if err != nil {
		return errs.Wrap(err)
	}

	fmt.Fprintf(f, "\tConfig[\"%s\"] = []ConfigKey{\n\t\t", name)
	if err != nil {
		return errs.Wrap(err)
	}

	err = writeConfigStruct(f, "STORJ", t)
	if err != nil {
		return errs.Wrap(err)
	}
	_, err = f.WriteString(`
   }
}
`)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func writeConfigStruct(f *os.File, prefix string, configType reflect.Type) error {
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			err := writeConfigStruct(f, prefix+"_"+camelToUpperCase(field.Name), field.Type)
			if err != nil {
				return errs.Wrap(err)
			}
		} else {
			name := prefix + "_" + camelToUpperCase(field.Name)
			_, err := fmt.Fprintf(f, `{
			Name:        "%s",
			Description: "%s",
			Default:     "%s",
		}, `, name, safe(field.Tag.Get("help")), safe(field.Tag.Get("default")))
			if err != nil {
				return errs.Wrap(err)
			}
		}
	}
	return nil
}

func safe(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

func camelToUpperCase(name string) string {
	smallCapital := regexp.MustCompile("([a-z])([A-Z])")
	name = smallCapital.ReplaceAllString(name, "${1}_$2")
	return strings.ToUpper(name)
}

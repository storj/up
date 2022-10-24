// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/zeebo/errs"

	"storj.io/gateway-mt/pkg/auth"
	"storj.io/gateway-mt/pkg/linksharing"
	"storj.io/storj-up/cmd/config"
	"storj.io/storj/satellite"
	"storj.io/storj/storagenode"
	"storj.io/storjscan"
)

var configTypes = map[string]reflect.Type{
	"storagenode":     reflect.TypeOf(storagenode.Config{}),
	"satellite-api":   reflect.TypeOf(satellite.Config{}),
	"satellite-admin": reflect.TypeOf(satellite.Config{}),
	"satellite-core":  reflect.TypeOf(satellite.Config{}),
	"linksharing":     reflect.TypeOf(linksharing.Config{}),
	"authservice":     reflect.TypeOf(auth.Config{}),
	"storjscan":       reflect.TypeOf(storjscan.Config{}),
}

func main() {
	templateDir := "."
	if len(os.Args) > 1 {
		templateDir = os.Args[1]
	}
	for name, t := range configTypes {
		err := generateSingle(templateDir, name, t)
		if err != nil {
			panic(err)
		}
	}
	err := generateCombiner(templateDir, configTypes)
	if err != nil {
		panic(err)
	}
}

func generateCombiner(templateDir string, types map[string]reflect.Type) error {
	fileName := "../all.go"
	fmt.Println("Writing " + fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		_ = f.Close()
	}()

	t, err := template.New("all.tpl").
		Funcs(map[string]interface{}{
			"goName": goName,
		}).
		ParseFiles(filepath.Join(templateDir, "all.tpl"))
	if err != nil {
		return errs.Wrap(err)
	}

	err = t.Execute(f, struct {
		Configs map[string]reflect.Type
	}{
		Configs: types,
	})
	return err
}

func generateSingle(templateDir string, name string, root reflect.Type) error {
	fileName := "../" + name + ".go"
	fmt.Println("Writing " + fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		_ = f.Close()
	}()

	t, err := template.New("single.tpl").
		Funcs(map[string]interface{}{
			"goName": goName,
		}).
		ParseFiles(filepath.Join(templateDir, "single.tpl"))
	if err != nil {
		return errs.Wrap(err)
	}

	options, err := collectOptions("STORJ", root)
	if err != nil {
		return errs.Wrap(err)
	}

	err = t.Execute(f, struct {
		Name    string
		Options []config.Option
	}{
		Name:    name,
		Options: options,
	})

	return err
}

func goName(name string) string {
	return strings.ReplaceAll(name, "-", "")
}

func collectOptions(prefix string, configType reflect.Type) (res []config.Option, err error) {
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			r, err := collectOptions(prefix+"_"+camelToUpperCase(field.Name), field.Type)
			if err != nil {
				return res, errs.Wrap(err)
			}
			res = append(res, r...)
		} else {
			name := prefix + "_" + camelToUpperCase(field.Name)
			res = append(res, config.Option{
				Name:        name,
				Description: safe(field.Tag.Get("help")),
				Default:     safe(field.Tag.Get("default"))})
		}
	}
	return res, nil
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

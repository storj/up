// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package runtime

import (
	"bytes"
	"text/template"

	"github.com/zeebo/errs/v2"
)

// Render can resolve all the go templates in a string.
func Render(r Runtime, serviceInstance ServiceInstance, original string) (string, error) {
	funcMap := map[string]interface{}{
		"Host": func(service string, hostType string) string {
			return r.GetHost(ServiceInstanceFromStr(service), hostType)
		},
		"Port": func(service string, portType string) int {
			return r.GetPort(ServiceInstanceFromStr(service), portType)
		},
		"Environment": func(service string, key string) (string, error) {
			val := r.Get(ServiceInstanceFromStr(service), key)
			if val == "" {
				return "", errs.Errorf("Variable is not specified in the environment '%s'/'%s'", service, key)
			}
			return val, nil
		},
	}

	tpl, err := template.New("line").
		Funcs(funcMap).
		Parse(original)
	if err != nil {
		return "", err
	}
	dest := bytes.Buffer{}
	err = tpl.Execute(&dest, struct {
		This string
	}{
		This: serviceInstance.String(),
	})
	if err != nil {
		return "", err
	}
	return dest.String(), nil

}

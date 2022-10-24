// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

func {{ .Name | goName }}Config() []Option {
	return []Option{
		{{ range .Options}}
		{
			Name:        "{{ .Name }}",
			Description: "{{ .Description }}",
			Default:     "{{ .Default }}",
		},{{ end }}
   }
}

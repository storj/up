// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

func All() map[string][]Option {
	return map[string][]Option{ {{ range $name, $def := .Configs }}
		"{{ $name }}": {{ $name | goName }}Config(),{{ end }}
	}
}

#!/usr/bin/env bash
cd $(dirname "${BASH_SOURCE[0]}")

set -euxo pipefail

{{ range $k,$v := .Service.Environment}}
export {{$k}}="{{$v}}"{{end}}
mkdir -p ./{{.Service.ID.Name}}/{{.Service.ID.Instance}}

#{{ if eq .Service.ID.Name "satellite-api"}}
#if [ ! -f ".{{.Service.ID.Name}}-migrated" ]; then
#    satellite run migration --identity-dir={{.Service.ID.Name}}/{{.Service.ID.Instance}} --config-dir={{.Service.ID.Name}}/{{.Service.ID.Instance}}
#    touch ".{{.Service.ID.Name}}-migrated"
#fi
#{{end}}

#RUN
{{range .Service.Command}}{{.}} {{end}}


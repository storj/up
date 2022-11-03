#!/usr/bin/env bash
cd $(dirname "${BASH_SOURCE[0]}")

set -euxo pipefail

{{ range $k,$v := .Service.Environment}}
export {{$k}}="{{$v}}"{{end}}
mkdir -p ./{{.Service.ID.Name}}/{{.Service.ID.Instance}}

#RUN
{{range .Service.Command}}{{.}} {{end}}


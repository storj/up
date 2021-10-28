package templates

import _ "embed"

//go:embed docker-compose.template.yaml
var ComposeTemplate []byte

//go:embed prometheus.yml
var PrometheusYaml []byte

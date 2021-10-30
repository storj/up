package dockerfiles

import _ "embed"

//go:embed storj.Dockerfile
var StorjDocker []byte

//go:embed edge.Dockerfile
var EdgeDocker []byte

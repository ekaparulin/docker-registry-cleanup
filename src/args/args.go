package args

import (
	"errors"
	"flag"
)

type Config struct {
	DockerRegistry   *string
	DockerRepository *string
	DryRun           *bool
	Verbose          *bool
	ListSchemaV1     *bool
}

func NewConfig() (Config, error) {
	c := Config{}

	c.DockerRegistry = flag.String("registry", "", "Docker registry URL")
	c.DockerRepository = flag.String("repo", "", "Docker image repository")
	c.DryRun = flag.Bool("dry-run", false, "Dry run")
	c.Verbose = flag.Bool("verbose", false, "Verbose mode")
	c.ListSchemaV1 = flag.Bool("list-schema-v1", false, "List Schema V1 images")

	flag.Parse()

	if *c.DockerRegistry == "" {
		return c, errors.New("DockerRegistry is not specified")
	}

	if *c.DockerRepository == "" {
		return c, errors.New("DockerRepository is not specified")
	}

	return c, nil
}

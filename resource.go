package streamline

import (
	"github.com/ory/dockertest"
)

// Resource defines a docker image along with the configuration
type Resource struct {
	// Wait is a function used to determine if a resource is ready
	// Returning an error will cause the resource to expontentially
	// block until it is ready.
	Wait func() error

	RunOpts  *dockertest.RunOptions
	resource *dockertest.Resource
}

// GetDockerResource returns the current resource
func (r *Resource) GetDockerResource() *dockertest.Resource {
	return r.resource
}

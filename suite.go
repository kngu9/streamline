package streamline

import (
	"fmt"
	"sync"

	jc "github.com/juju/testing/checkers"
	"github.com/ory/dockertest"
	gc "gopkg.in/check.v1"
)

// Suite defines the interface for streamline's suite
type Suite interface {
	AddResource(string, Resource) error
	DeleteResource(name string) error
	AddDefaultResource(name string, rType ResourceType) error
	GetResource(name string) (*Resource, error)

	SetUpTest(*gc.C)
	TearDownTest(*gc.C)
}

// Suite is the main streamline gocheck suite
type suite struct {
	pool        *dockertest.Pool
	resourceMap map[string]*Resource

	testingLock sync.Mutex
}

// New creates a new streamline test suite
func New(address string) (Suite, error) {
	s := suite{
		resourceMap: make(map[string]*Resource),
	}

	pool, err := dockertest.NewPool(address)
	if err != nil {
		return nil, err
	}
	s.pool = pool

	return &s, nil
}

// AddResource is a synchronous operation that blocks if there's a test running.
// The reason behind this is because we do not want to spin up more resources,
// while there are resources currently spinning.
func (s *suite) AddResource(name string, resource Resource) error {
	s.testingLock.Lock()
	defer s.testingLock.Unlock()

	if _, ok := s.resourceMap[name]; ok {
		return fmt.Errorf("resource name %s already exists, please remove it or use a different name", name)
	}

	s.resourceMap[name] = &resource

	return nil
}

func (s *suite) AddDefaultResource(name string, rType ResourceType) error {
	var resource *Resource

	switch rType {
	case ResourceTypePostgres:
		resource = &Resource{
			RunOpts: pgResource,
			Wait:    pgWait(s, name),
		}
	default:
		return fmt.Errorf("invalid resource type %+v", rType)
	}

	return s.AddResource(name, *resource)
}

// DeleteResource is a synchronous operation that blocks if there's a test running.
func (s *suite) DeleteResource(name string) error {
	s.testingLock.Lock()
	defer s.testingLock.Unlock()

	if _, ok := s.resourceMap[name]; ok {
		delete(s.resourceMap, name)
		return nil
	}
	return fmt.Errorf("resource %s not found", name)
}

// GetResource returns the resource attributed to the name specified
func (s *suite) GetResource(name string) (*Resource, error) {
	if val, ok := s.resourceMap[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("resource %s does not exist", name)
}

func (s *suite) SetUpTest(c *gc.C) {
	s.testingLock.Lock()

	for k, res := range s.resourceMap {
		resource, err := s.pool.RunWithOptions(res.RunOpts)
		c.Assert(err, jc.ErrorIsNil)
		s.resourceMap[k].resource = resource
		fmt.Println(s.resourceMap[k])
		if res.Wait != nil {
			c.Assert(s.pool.Retry(res.Wait), jc.ErrorIsNil)
		}
	}
}

func (s *suite) TearDownTest(c *gc.C) {
	defer s.testingLock.Unlock()

	for _, res := range s.resourceMap {
		c.Assert(res.resource.Close(), jc.ErrorIsNil)
	}
}

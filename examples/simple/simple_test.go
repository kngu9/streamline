package simple_test

import (
	"context"
	"database/sql"
	"fmt"
	stdtesting "testing"

	jc "github.com/juju/testing/checkers"
	_ "github.com/lib/pq" // default pg driver
	gc "gopkg.in/check.v1"

	"github.com/kngu9/streamline"
)

func Test(t *stdtesting.T) { gc.TestingT(t) }

var _ = gc.Suite(&SimpleSuite{})

type SimpleSuite struct {
	streamline.Suite

	db *sql.DB
}

func (s *SimpleSuite) SetUpSuite(c *gc.C) {
	suite, err := streamline.New("")
	c.Assert(err, jc.ErrorIsNil)
	s.Suite = suite

	s.Suite.AddDefaultResource("db", streamline.ResourceTypePostgres)
}

func (s *SimpleSuite) SetUpTest(c *gc.C) {
	s.Suite.SetUpTest(c)

	res, err := s.Suite.GetResource("db")
	c.Assert(err, jc.ErrorIsNil)
	c.Log(res)

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://test:test@localhost:%s/test?sslmode=disable", res.GetDockerResource().GetPort("5432/tcp")))
	c.Assert(err, jc.ErrorIsNil)
	s.db = db
}

func (s *SimpleSuite) TestExample(c *gc.C) {
	row := s.db.QueryRowContext(context.Background(), `SELECT 'hello world'`)
	c.Assert(row, gc.NotNil)

	var msg string
	c.Assert(row.Scan(&msg), jc.ErrorIsNil)
	c.Assert(msg, gc.Equals, "hello world")
}

func (s *SimpleSuite) TearDownTest(c *gc.C) {
	if s.db != nil {
		c.Assert(s.db.Close(), jc.ErrorIsNil)
	}

	s.Suite.TearDownTest(c)
}

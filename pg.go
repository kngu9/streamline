package streamline

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // default pg driver
	"github.com/ory/dockertest"
)

var pgResource = &dockertest.RunOptions{
	Repository: "postgres",
	Tag:        "latest",
	Env: []string{
		"POSTGRES_USER=test",
		"POSTGRES_PASSWORD=test",
		"POSTGRES_DB=test",
	},
}

func pgWait(suite *suite, name string) func() error {
	return func() error {
		curRes := suite.resourceMap[name]
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://test:test@localhost:%s/test?sslmode=disable", curRes.resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}
}

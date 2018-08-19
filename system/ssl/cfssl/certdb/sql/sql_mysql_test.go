// +build mysql

package sql

import (
	"testing"

	"github.com/ironstar-io/tokaido/system/ssl/cfssl/certdb/testdb"
)

func TestMySQL(t *testing.T) {
	db := testdb.MySQLDB()
	ta := TestAccessor{
		Accessor: NewAccessor(db),
		DB:       db,
	}
	testEverything(ta, t)
}

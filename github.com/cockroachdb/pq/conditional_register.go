// +build register_crdb

// andrei: Unconditional registration of the pq driver under the name "postgres"
// was remove from conn.go and added to this conditionally-built file in order
// to support linking this forked lib/pq with an unforked one. Users of the
// forked lib don't need to use it through Go's sql package, so they don't need
// this registration. pq's own tests do, however, so they must be built with
// this build tag.

package pq

import "database/sql"

func init() {
	sql.Register("postgres", &drv{})
}

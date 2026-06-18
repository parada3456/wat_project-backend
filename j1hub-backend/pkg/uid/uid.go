package uid

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

func New(prefix string) string {
	return prefix + ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

// prefixes: usr_ prf_ frn_ phs_ uph_ mis_ ums_ tsk_ utk_
//           ldg_ bdg_ ubdg_ crd_ txn_ spl_ job_ hsg_ crt_ smr_ rvw_

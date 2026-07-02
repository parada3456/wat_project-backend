package uid_test

import (
	"strings"
	"testing"

	"github.com/parada3456/wat_project-backend/pkg/uid"
	"github.com/stretchr/testify/assert"
)

func TestNew_Prefix(t *testing.T) {
	id := uid.New("usr_")
	assert.True(t, strings.HasPrefix(id, "usr_"))
	assert.Len(t, id, 4+26) // "usr_" (4 chars) + ULID (26 chars)
}

package channelz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Just a simple smoke test
func TestCreateHandler(t *testing.T) {
	h := CreateHandler("/prefix", ":8080")
	assert.NotNil(t, h)
}

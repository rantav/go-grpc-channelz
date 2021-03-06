package channelz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteServerPage(t *testing.T) {
	assert := assert.New(t)
	handler := grpcChannelzHandler{client: &mockChannelzClient{}}
	var b strings.Builder
	handler.WriteServerPage(&b, 6)
	assert.Contains(b.String(), `ChannelZ server 6`)
	assert.Contains(b.String(), `CT_INFO [1970-01-01T00:00:06Z]: setup`)
}

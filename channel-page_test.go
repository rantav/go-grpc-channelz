package channelz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteChannelPage(t *testing.T) {
	assert := assert.New(t)
	handler := grpcChannelzHandler{client: &mockChannelzClient{}}
	var b strings.Builder
	handler.WriteChannelPage(&b, 2)
	assert.Contains(b.String(), "channel 2")
	assert.Contains(b.String(), "CT_INFO [1970-01-01T00:00:06Z]: setup")
}

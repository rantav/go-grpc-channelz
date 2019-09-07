package channelz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteSocketPage(t *testing.T) {
	assert := assert.New(t)
	handler := grpcChannelzHandler{client: &mockChannelzClient{}}
	var b strings.Builder
	handler.WriteSocketPage(&b, 9)
	assert.Contains(b.String(), `ChannelZ socket 9`)
	assert.Contains(b.String(), `hello: world`)
}

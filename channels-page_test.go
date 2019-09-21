package channelz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteChannelsPage(t *testing.T) {
	assert := assert.New(t)
	handler := grpcChannelzHandler{client: &mockChannelzClient{}}
	var b strings.Builder
	handler.WriteChannelsPage(&b, 2)
	assert.Contains(b.String(), `<a href="/subchannel/8"><b>8</b> eight</a>`)
}

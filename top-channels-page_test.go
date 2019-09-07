package channelz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteTopChannelsPage(t *testing.T) {
	assert := assert.New(t)
	handler := grpcChannelzHandler{client: &mockChannelzClient{}}
	var b strings.Builder
	handler.WriteTopChannelsPage(&b)
	assert.Contains(b.String(), `<a href="/channelz/subchannel/8"><b>8</b> eight</a>`)
	assert.Contains(b.String(), `<a href="/channelz/server/1"><b>1</b> one</a>`)
}

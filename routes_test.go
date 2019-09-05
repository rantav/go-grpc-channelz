package channelz

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
}

func (m *mockHandler) WriteTopChannelsPage(w io.Writer) {
	w.Write([]byte("top"))
}
func (m *mockHandler) WriteChannelPage(w io.Writer, c int64) {
	w.Write([]byte(fmt.Sprintf("channel %d", c)))
}
func (m *mockHandler) WriteSubchannelPage(w io.Writer, c int64) {
	w.Write([]byte(fmt.Sprintf("subchannel %d", c)))
}
func (m *mockHandler) WriteServerPage(w io.Writer, c int64) {
	w.Write([]byte(fmt.Sprintf("server %d", c)))
}

func TestCreateRouter(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	r := createRouter("/channelz", &mockHandler{})
	assert.NotNil(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	expects := map[string]string{
		"/channelz":              "top",
		"/channelz/channel/4":    "channel 4",
		"/channelz/subchannel/5": "subchannel 5",
		"/channelz/server/3":     "server 3",
	}
	for route, expected := range expects {
		res, err := http.Get(ts.URL + route)
		require.NoError(err)

		response, err := ioutil.ReadAll(res.Body)
		require.NoError(err)
		res.Body.Close()

		assert.Equal([]byte(expected), response, "For path %s the expected result was %q, but instead we got %q",
			route, expected, response)
	}

}

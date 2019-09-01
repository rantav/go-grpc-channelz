package channelz

import "io"

func write404(w io.Writer) {
	writeHeader(w, "ChannelZ 404 not found")
	writeFooter(w)
}

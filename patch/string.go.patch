package spdy

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

///////////////////
// FIXME!!
///////////////////

// join slice element with separator
// arrayJoin([a, b, c], ",") => "a, b, c"
func join(arr []string, sep string) string {
	var str string
	length := len(arr) - 1
	for i, v := range arr {
		str += v
		if i == length {
			break
		}
		str += sep
	}
	return str
}

func formatHeader(headers http.Header) []string {
	arr := make([]string, len(headers))
	for k, v := range headers {
		k = strings.ToLower(k)
		switch k {
		case "host", "method", "path", "scheme", "version":
			k = ":"
		case "url":
			k = ":path"
		}
		arr = append(arr, fmt.Sprintf("%s: %s\n", k, join(v, ",")))
	}
	sort.Strings(arr)
	return arr
}

func (f *SynStreamFrame) String() string {
	var str string
	str += fmt.Sprintf("SYN_STREAM frame <version=%d, flags=%d, length=%d>\n", f.CFHeader.version, f.CFHeader.Flags, f.CFHeader.length)

	str += fmt.Sprintf("(stream_id=%d, assoc_stream_id=%d, pri=%d)\n", f.StreamId, f.AssociatedToStreamId, f.Priority)
	headers := formatHeader(f.Headers)

	for _, v := range headers {
		str += fmt.Sprintf("%s", v)
	}

	return str
}

func (f *SynReplyFrame) String() string {
	var str string
	str += fmt.Sprintf("SYN_REPLY frame <version=%d, flags=%d, length=%d>\n", f.CFHeader.version, f.CFHeader.Flags, f.CFHeader.length)

	str += fmt.Sprintf("(stream_id=%d)\n", f.StreamId)
	headers := formatHeader(f.Headers)

	for _, v := range headers {
		str += fmt.Sprintf("%s", v)
	}

	return str
}

func (f *DataFrame) String() string {
	var str string
	str += fmt.Sprintf("DATA frame (stream_id=%d, flags=%d, length=%d)", f.StreamId, f.Flags, len(f.Data))
	return str
}

func (f *GoAwayFrame) String() string {
	var str string
	str += fmt.Sprintf("GOAWAY frame <status=%d>\n", f.Status)
	str += fmt.Sprintf("(last_good_stream_id=%d)\n", f.LastGoodStreamId)
	return str
}

func (f *SettingsFrame) String() string {
	var str string
	str += "SETTINGS frame\n"
	str += fmt.Sprintf("(niv=%d)\n", len(f.FlagIdValues))
	for _, fiv := range f.FlagIdValues {
		str += fmt.Sprintf("[%d(%d):%d]\n", fiv.Id, fiv.Flag, fiv.Value)
	}
	return str
}

func (f *RstStreamFrame) String() string {
	return "RstStreamFrame"
}

func (f *WindowUpdateFrame) String() string {
	return "WindowUpdateFrame"
}

func (f *PingFrame) String() string {
	return "PingFrame"
}

func (f *HeadersFrame) String() string {
	return "HeadersFrame"
}

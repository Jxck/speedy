package speedy

import (
	"code.google.com/p/go.net/spdy"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

// debug data
var ResponseHtml = `
<html>
	<head>
		<title>SPDY</title>
		<script type="text/javascript" src="https://localhost:3000/test.js"></script>
	</head>
	<body>
		<h1>Speedy :)</h1>
	</body>
</html>
`

var HeadersFixtureHtml = http.Header{
	":version":       []string{"http/1.1"},
	":status":        []string{"200 OK"},
	":host":          []string{"localhost:3000"},
	":path":          []string{"/"},
	":scheme":        []string{"https"},
	"location":       []string{"https://localhost:3000/"},
	"content-type":   []string{"text/html; charset=utf-8"},
	"content-length": []string{strconv.Itoa(len(ResponseHtml))},
	"server":         []string{"speedy"},
}

var ResponseJS = `
console.log("Speedy");
`

var HeadersFixtureJS = http.Header{
	":version":       []string{"http/1.1"},
	":status":        []string{"200 OK"},
	":host":          []string{"localhost:3000"},
	":path":          []string{"/test.js"},
	":scheme":        []string{"https"},
	"location":       []string{"https://localhost:3000/test.js"},
	"content-type":   []string{"text/javascript; charset=utf-8"},
	"content-length": []string{strconv.Itoa(len(ResponseJS))},
	"server":         []string{"speedy"},
}

func sendSynReply(header http.Header, framer *spdy.Framer, frame *spdy.SynStreamFrame) error {

	synReply := spdy.SynReplyFrame{
		CFHeader: spdy.ControlFrameHeader{}, // Flag is 0x00
		StreamId: frame.StreamId,
		Headers:  header,
	}

	err := framer.WriteFrame(&synReply)
	if err != nil {
		return err
	}
	debug("send %v", &synReply)
	return nil
}

func sendData(data string, framer *spdy.Framer, frame *spdy.SynStreamFrame) error {
	dataFrame := spdy.DataFrame{
		StreamId: frame.StreamId,
		Flags:    spdy.DataFlagFin,
		Data:     []byte(data),
	}

	err := framer.WriteFrame(&dataFrame)
	if err != nil {
		return err
	}
	debug("send %v", &dataFrame)
	return nil
}

func sendSynStream(header http.Header, framer *spdy.Framer, frame *spdy.SynStreamFrame) error {
	synStreamFrame := spdy.SynStreamFrame{
		CFHeader: spdy.ControlFrameHeader{
			Flags: spdy.ControlFlagUnidirectional,
		},
		StreamId:             frame.StreamId + 1,
		Headers:              header,
		AssociatedToStreamId: frame.StreamId,
	}

	err := framer.WriteFrame(&synStreamFrame)
	if err != nil {
		return err
	}

	debug("push %v", &synStreamFrame)
	return nil
}

func pushData(data string, framer *spdy.Framer, frame *spdy.SynStreamFrame) error {
	dataFrame := spdy.DataFrame{
		StreamId: frame.StreamId + 1,
		Flags:    spdy.DataFlagFin,
		Data:     []byte(data),
	}

	err := framer.WriteFrame(&dataFrame)
	if err != nil {
		return err
	}
	debug("push %v", &dataFrame)
	return nil
}

type HandleHello struct{}

func (h HandleHello) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello")
}

type ResponseWriter struct{}

func (r ResponseWriter) Header() http.Header {
	return http.Header{}
}

func (r ResponseWriter) Write(b []byte) (int, error) {
	return 1, nil
}

func (r ResponseWriter) WriteHeader(i int) {
}

func handleSynStreamFrame(framer *spdy.Framer, synStream *spdy.SynStreamFrame) error {
	debug("recv %v", synStream.Headers)

	method := synStream.Headers.Get(":method")
	path := synStream.Headers.Get(":path")
	pipeReader, _ := io.Pipe()
	req, err := http.NewRequest(method, path, pipeReader)
	if err != nil {
		return err
	}

	res := ResponseWriter{}
	debug("%v", req)
	debug("%v", res)

	// send reply
	err = sendSynStream(HeadersFixtureJS, framer, synStream)
	if err != nil {
		return err
	}

	// send data
	err = pushData(ResponseJS, framer, synStream)
	if err != nil {
		return err
	}

	// send reply
	err = sendSynReply(HeadersFixtureHtml, framer, synStream)
	if err != nil {
		return err
	}

	// send data
	err = sendData(ResponseHtml, framer, synStream)
	if err != nil {
		return err
	}

	return nil
}

func handleGoAwayFrame(framer *spdy.Framer, goAway *spdy.GoAwayFrame) error {
	debug("recv %v", goAway)
	return nil
}

func handleConnection(conn net.Conn) error {
	framer, err := spdy.NewFramer(conn, conn)
	if err != nil {
		return err
	}

	for {
		frame, err := framer.ReadFrame()
		if err != nil {
			return err
		}

		switch frame := frame.(type) { // decide type
		case *spdy.SynStreamFrame:
			err := handleSynStreamFrame(framer, frame)
			if err != nil {
				return err
			}
		case *spdy.GoAwayFrame:
			err := handleGoAwayFrame(framer, frame)
			if err != nil {
				return err
			}
		case *spdy.SynReplyFrame:
			debug("recv %v", frame)
		case *spdy.RstStreamFrame:
			debug("recv %v", frame)
		case *spdy.SettingsFrame:
			debug("recv %v", frame)
		case *spdy.WindowUpdateFrame:
			debug("recv %v", frame)
		case *spdy.PingFrame:
			debug("recv %v", frame)
		case *spdy.HeadersFrame:
			debug("recv %v", frame)
		case *spdy.DataFrame:
			debug("recv %v", frame)
		default:
			debug("unknown frame")
			log.Fatalf("unknown frame %v", frame)
		}
	}
	return nil
}

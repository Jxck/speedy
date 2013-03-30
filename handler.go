package speedy

import (
	"github.com/jxck/speedy/spdy"
	"io"
	"log"
	"net"
	"net/http"
)

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

	method := synStream.Headers.Get(":method")
	path := synStream.Headers.Get(":path")
	pipeReader, _ := io.Pipe()
	req, err := http.NewRequest(method, path, pipeReader)
	if err != nil {
		return err
	}

	res := ResponseWriter{}
	_, _ = req, res

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
	defer conn.Close()
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

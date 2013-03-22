package speedy

import (
	"code.google.com/p/go.net/spdy"
	"log"
	"net"
	"net/http"
)

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

		switch frametype := frame.(type) {
		case *spdy.SynStreamFrame:
			synStream := frametype //frame.(*spdy.SynStreamFrame)
			debug("recv %v", synStream)

			// debug data
			var ResponseData = "<h1>hello</h1>"

			var HeadersFixture = http.Header{
				":version":       []string{"http/1.1"},
				":status":        []string{"200 OK"},
				"location":       []string{"http://localhost:3000/"},
				"content-type":   []string{"text/html; charset=utf-8"},
				"content-length": []string{"14"},
				"server":         []string{"speedy"},
			}

			// send reply
			synReply := spdy.SynReplyFrame{
				CFHeader: spdy.ControlFrameHeader{}, //Flag is 0x00
				StreamId: synStream.StreamId,
				Headers:  HeadersFixture,
			}

			err := framer.WriteFrame(&synReply)
			if err != nil {
				return err
			}
			debug("send %v", &synReply)

			// send data
			dataFrame := spdy.DataFrame{
				StreamId: synStream.StreamId,
				Flags:    spdy.DataFlagFin,
				Data:     []byte(ResponseData),
			}

			err = framer.WriteFrame(&dataFrame)
			if err != nil {
				return err
			}
			debug("send %v", &dataFrame)

		default:
			debug("unknown frame")
			log.Fatalf("unknown frame %v", frametype)
		}
	}
	return nil
}

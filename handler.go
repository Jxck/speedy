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
			synframe := frame.(*spdy.SynStreamFrame)
			debug("%v", synframe)

			// Debug data
			var HeadersFixture = http.Header{
				"Url":     []string{"http://localhost:3000/"},
				"Method":  []string{"get"},
				"Version": []string{"http/1.1"},
			}

			// send reply
			synReplyFrame := spdy.SynReplyFrame{
				CFHeader: spdy.ControlFrameHeader{}, //Flag is 0x00
				StreamId: synframe.StreamId,
				Headers:  HeadersFixture,
			}

			err := framer.WriteFrame(&synReplyFrame)
			if err != nil {
				return err
			}

			// send data
			dataFrame := spdy.DataFrame{
				StreamId: synframe.StreamId,
				Flags:    spdy.DataFlagFin,
				Data:     []byte{'h', 'e', 'l', 'l', 'o'},
			}

			err = framer.WriteFrame(&dataFrame)
			if err != nil {
				return err
			}

		default:
			debug("unknown frame %v", frametype)
			log.Fatalf("unknown frame %v", frametype)
		}
	}
	return nil
}

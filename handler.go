package speedy

import (
	"code.google.com/p/go.net/spdy"
	"log"
	"net"
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

		default:
			log.Fatalf("unknown frame %v", frametype)
		}
	}
	return nil
}

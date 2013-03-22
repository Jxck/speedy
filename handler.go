package speedy

import (
	"code.google.com/p/go.net/spdy"
	"log"
	"net"
	"net/http"
	"strconv"
)

// debug data
var ResponseData = `
<html>
	<head>
		<title>SPDY</title>
		<script type="text/javascript" src="test.js"></script>
	</head>
	<body>
		<h1>Speedy :)</h1>
	</body>
</html>
`

func handleSynStreamFrame(framer *spdy.Framer, synStream *spdy.SynStreamFrame) error {
	debug("recv %v", synStream)

	var HeadersFixture = http.Header{
		":version":       []string{"http/1.1"},
		":status":        []string{"200 OK"},
		"location":       []string{"http://localhost:3000/"},
		"content-type":   []string{"text/html; charset=utf-8"},
		"content-length": []string{strconv.Itoa(len(ResponseData))},
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

		switch frametype := frame.(type) {
		case *spdy.SynStreamFrame:
			err := handleSynStreamFrame(framer, frametype)
			if err != nil {
				return err
			}
		case *spdy.GoAwayFrame:
			err := handleGoAwayFrame(framer, frametype)
			if err != nil {
				return err
			}
		case *spdy.SynReplyFrame:
			debug("recv %v", frametype)
		case *spdy.RstStreamFrame:
			debug("recv %v", frametype)
		case *spdy.SettingsFrame:
			debug("recv %v", frametype)
		case *spdy.WindowUpdateFrame:
			debug("recv %v", frametype)
		case *spdy.PingFrame:
			debug("recv %v", frametype)
		case *spdy.HeadersFrame:
			debug("recv %v", frametype)
		case *spdy.DataFrame:
			debug("recv %v", frametype)
		default:
			debug("unknown frame")
			log.Fatalf("unknown frame %v", frametype)
		}
	}
	return nil
}

package speedy

import (
	"crypto/tls"
	"net"
	"net/http"
)

// make tls.Config with spdy/3 npn
func ConfigTLS(addr, certFile, keyFile string) (*tls.Config, error) {
	if addr == "" {
		addr = ":https"
	}
	var config *tls.Config = &tls.Config{}
	config.NextProtos = []string{"spdy/3"}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ListenAndServe(addr, certFile, keyFile string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// configure tls with spdy/3 npn
	config, err := ConfigTLS(addr, certFile, keyFile)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(listener, config)
	debug(clr.Cyan("listening server %s"), tlsListener.Addr())

	return nil
}

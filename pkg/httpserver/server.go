package httpserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	cfg             *config.Config
}

// New -.
func New(handler http.Handler, cfg *config.Config, opts ...Option) *Server {

	httpServer := &http.Server{}

	if (cfg.App.Environment == "development" || cfg.App.Environment == "staging") && cfg.HTTPServer.UseSSL {
		cfgClient := &tls.Config{
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		httpServer.Addr = ":" + cfg.HTTPServer.Port
		httpServer.Handler = handler
		httpServer.ReadTimeout = _defaultReadTimeout
		httpServer.WriteTimeout = _defaultWriteTimeout
		httpServer.TLSConfig = cfgClient
		httpServer.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	} else {
		httpServer.Handler = handler
		httpServer.ReadTimeout = _defaultReadTimeout
		httpServer.WriteTimeout = _defaultWriteTimeout
		httpServer.Addr = _defaultAddr
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		cfg:             cfg,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		if s.cfg.HTTPServer.UseSSL {
			s.notify <- s.server.ListenAndServeTLS(s.cfg.BaseDir+s.cfg.HTTPServer.SSLCert, s.cfg.BaseDir+s.cfg.HTTPServer.SSLKey)
		} else {
			s.notify <- s.server.ListenAndServe()

		}
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

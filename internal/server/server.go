package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/maxrasky/faraway/internal/utils"
)

type Quoter interface {
	GetQuote() string
}

type Verifier interface {
	Challenge() []byte
	Verify(challenge, solution []byte) error
}

type Server struct {
	listener net.Listener
	cfg      *Config
	verifier Verifier
	quotes   Quoter
}

func NewServer(cfg *Config, verifier Verifier, quoter Quoter) *Server {
	return &Server{
		cfg:      cfg,
		verifier: verifier,
		quotes:   quoter,
	}
}

func (s *Server) Run(ctx context.Context) (err error) {
	s.listener, err = net.Listen("tcp", s.cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Info().Msg("started server")

	go func() {
		<-ctx.Done()
		_ = s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			continue
		}

		go func(conn net.Conn) {
			if err = s.handle(conn); err != nil {
				log.Err(err).Msg("failed to handle connection")
			}
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) error {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(s.cfg.Deadline))

	if _, err := utils.ReadMessage(conn); err != nil {
		return fmt.Errorf("read message err: %w", err)
	}

	challenge := s.verifier.Challenge()
	if err := utils.WriteMessage(conn, challenge); err != nil {
		return fmt.Errorf("send challenge err: %w", err)
	}

	solution, err := utils.ReadMessage(conn)
	if err != nil {
		return fmt.Errorf("receive proof err: %w", err)
	}

	if err = s.verifier.Verify(challenge, solution); err != nil {
		return fmt.Errorf("invalid verify: %w", err)
	}

	quote := s.quotes.GetQuote()
	if err = utils.WriteMessage(conn, []byte(quote)); err != nil {
		return fmt.Errorf("send quote err: %w", err)
	}

	return nil
}

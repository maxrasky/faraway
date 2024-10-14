package client

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/maxrasky/faraway/internal/utils"
)

type Solver interface {
	Solve(challenge []byte) ([]byte, error)
}

type Client struct {
	solver     Solver
	serverAddr string
}

func NewClient(serverAddr string, solver Solver) *Client {
	return &Client{
		solver:     solver,
		serverAddr: serverAddr,
	}
}

func (c *Client) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		quote, err := c.GetQuote(ctx)
		if err != nil {
			log.Err(err).Msg("failed to get quote")
		} else {
			log.Info().Bytes("msg", quote).Msg("got a quote")
		}
	}
}

func (c *Client) GetQuote(ctx context.Context) ([]byte, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", c.serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Err(err).Msg("failed to close connection")
		}
	}()

	if err := utils.WriteMessage(conn, []byte("challenge")); err != nil {
		return nil, fmt.Errorf("send challenge request err: %w", err)
	}

	challenge, err := utils.ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("receive challenge err: %w", err)
	}

	solution, err := c.solver.Solve(challenge)
	if err != nil {
		return nil, fmt.Errorf("failed to solve challenge: %w", err)
	}

	if err := utils.WriteMessage(conn, solution); err != nil {
		return nil, fmt.Errorf("send solution err: %w", err)
	}

	quote, err := utils.ReadMessage(conn)
	if err != nil {
		return nil, fmt.Errorf("receive quote err: %w", err)
	}

	return quote, nil
}

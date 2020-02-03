package server

import "context"

// Server represents a patcher server
type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new server instance
func New(ctx context.Context, cancel context.CancelFunc) (*Server, error) {
	s := new(Server)
	s.ctx, s.cancel = ctx, cancel
	return s, nil
}

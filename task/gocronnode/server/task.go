package server

import "context"

type TaskFunc func(ctx context.Context, params ...string) (string, error)

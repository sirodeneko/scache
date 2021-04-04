package scache

import "context"

type GetF func(context.Context) (interface{}, error)

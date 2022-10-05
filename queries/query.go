package cqrs_queries

import "context"

type IQuery interface {
}

type IResponse interface {
}

type IQueryHandler interface {
	HandleQuery(ctx context.Context, query IQuery) (IResponse, error)
}

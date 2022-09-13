package cqrs_queries

type IQuery interface {
}

type IResponse interface {
}

type IQueryHandler interface {
	HandleQuery(query IQuery) (IResponse, error)
}

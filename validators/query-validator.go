package cqrs_validators

import queries "github.com/mitz-it/golang-cqrs/queries"

type IQueryValidator interface {
	Validate(query queries.IQuery) error
}

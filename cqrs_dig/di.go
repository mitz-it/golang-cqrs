package cqrs_dig

import (
	cqrs "github.com/mitz-it/golang-cqrs"
	"go.uber.org/dig"
)

func ProvideCommandHandler[TCommand any, TResponse any](container *dig.Container, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(handler cqrs.ICommandHandler[TCommand, TResponse]) error {
		return cqrs.RegisterCommandHandler(handler)
	})

	return err
}

func ProvideQueryHandler[TQuery any, TResponse any](container *dig.Container, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(handler cqrs.IQueryHandler[TQuery, TResponse]) error {
		return cqrs.RegisterQueryHandler(handler)
	})

	return err
}

func ProvideCommandBehavior[TBehavior cqrs.IBehavior](container *dig.Container, order int, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(behavior TBehavior) error {
		return cqrs.RegisterCommandBehavior(order, behavior)
	})

	return err
}

func ProvideQueryBehavior[TBehavior cqrs.IBehavior](container *dig.Container, order int, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(behavior TBehavior) error {
		return cqrs.RegisterQueryBehavior(order, behavior)
	})

	return err
}

func ProvideEventHandler[TEvent any](container *dig.Container, constructor interface{}) error {
	err := container.Provide(constructor)

	if err != nil {
		return err
	}

	err = container.Invoke(func(handler cqrs.IEvenHandler[TEvent]) error {
		return cqrs.RegisterEventSubcriber(handler)
	})

	return err
}

# Go - CQRS

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mitz-it_golang-cqrs&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=mitz-it_golang-cqrs) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=mitz-it_golang-cqrs&metric=coverage)](https://sonarcloud.io/summary/new_code?id=mitz-it_golang-cqrs)

A mediator pattern abstraction destined for CQRS usage.

Strongly inspired by:

- [MediatR](https://github.com/jbogard/MediatR)
- [Go-MediatR](https://github.com/mehdihadeli/Go-MediatR)

## Installation

```bash
go get -u github.com/mitz-it/golang-cqrs
```

## Commands Usage

```go
// Define a command.
type CreateProduct struct {
  // ...
}

// Define a response.
type Product struct {
  // ...
}

// Define a handler.
type CreateProductHandler struct {
  // ...
}

// Implement the ICommandHandler interface.
func (h *CreateProductHandler) Handle(ctx context.Context, c *CreateProduct) (*Product, error) {
  // ...
}

// Register the handler.
handler := &CreateProductHandler{}
cqrs.RegisterCommandHandler[*CreateProduct, *Product](handler)

// Send the command.
command := &CreateProduct {
  // ...
}

ctx := context.Background() // When using with OpenTelemetry, be sure to use the received context to propagate it.
product, err := cqrs.Send[*CreateProduct, *Product](ctx, command)
```

## Queries Usage

```go
// Define a query.
type GetProduct struct {
  // ...
}

// Define a response.
type Product struct {
  // ...
}

// Define a handler.
type GetProductHandler struct {
  // ...
}

// Implement the IQueryHandler interface.
func (h *GetProductHandler) Handle(ctx context.Context, c *GetProduct) (*Product, error) {
  // ...
}

// Register the handler.
handler := &GetProductHandler{}
cqrs.RegisterQueryHandler[*GetProduct, *Product](handler)

// Request the query data.
query := &GetProduct {
  // ...
}

ctx := context.Background() // When using with OpenTelemetry, be sure to use the received context to propagate it.
product, err := cqrs.Request[*GetProduct, *Product](ctx, query)
```

## Behaviors Usage

Behaviors can be shared between commands and queries, but they need to be registered separately.

When registering behaviors, you must set the priority order to execute behaviors.
Priority 0 will be executed first, then 1, and so on.

```go
// Define the behavior
type LoggingBehavior struct {
  // ...
}

// Implement the IBehavior Interface
func (b *LoggingBehavior) Handle(ctx context.Context, request interface{}, next cqrs.NextFunc) (interface{}, error) {
  logger.Log.Info("processing request...")
  res, err := next()
  logger.Log.Info("request processed...")
  return res, err
}

// Register the behavior for commands
behavior := &LoggingBehavior{}
order := 0
cqrs.RegisterCommandBehavior(order, behavior)

// Register the behavior for queries
behavior := &LoggingBehavior{}
order := 0
cqrs.RegisterQueryBehavior(order, behavior)
```

## Events Usage

```go
// Create the event
type ProductCreated struct {
  // ...
}

// Define the event handler
type ProductCreatedHandler struct {
  // ...
}

// Implement the IEventHandler interface
func (h *ProductCreatedHandler) Handle(ctx context.Context, event *ProductCreated) error {
  // ...
}

// Register the event handler
handler := &ProductCreatedHandler{}
cqrs.RegisterEventSubscriber[*ProductCreated](handler)


// Send an event synchronously
event := &ProductCreated{}
ctx := context.Background() // When using with OpenTelemetry, be sure to use the received context to propagate it.
err := cqrs.PublisEvent(ctx, event)

// Send a fire and forget event
cqrs.Listen() // Call this method only once in your application, like at main.go
ctx := context.Background() // When using with OpenTelemetry, be sure to use the received context to propagate it.
cqrs.PublisEventAsync(ctx, event)
```

## Domain Events Usage

Use the [Events Usage](#events-usage) as the setup for this example.

```go
// implement the INotifiable interface
type Product struct {
  // ...
  events []interface{}
}

func (p *Product) AddEvent(event interface{}) {
  p.events = append(p.events, event)
}

func (p *Product) ClearEvents() {
  p.events = []interface{}{}
}

func (p *Product) GetEvents() []interface{} {
  return p.events
}


// Create a behavior to send events
type DispatchEventsBehavior struct {
}

func (behavior *DispatchEventsBehavior) Handle(ctx context.Context, request interface{}, next cqrs.NextFunc) (interface{}, error) {
  res, err := next()

  if err != nil {
    return res, err
  }

  n, ok := res.(cqrs.INotifiable)

  if !ok {
    return res, err
  }

  for _, event := range n.GetEvents() {
    cqrs.PublishEventAsync(ctx, event)
  }

  return res, errs
}
```
# Go - CQRS

A mediator pattern abstraction destined for CQRS usage.

Strongly inspired by:

- [MediatR](https://github.com/jbogard/MediatR)
- [Go-MediatR](https://github.com/mehdihadeli/Go-MediatR)

## Installation

```bash
go get github.com/mitz-it/golang-cqrs@v1.2.0
```

## Usage

### Commands

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
mediator.RegisterCommandHandler[*CreateProduct, *Product](handler)

// Send the command.
command := &CreateProduct {
  // ...
}

ctx := context.Background() // When using with OpenTelemetry, be sure to use the received context to propagate it.
product, err := mediator.Send[*CreateProduct, *Product](ctx, command)
```

### Queries

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
mediator.RegisterQueryHandler[*GetProduct, *Product](handler)

// Request the query data.
query := &GetProduct {
  // ...
}

ctx := context.Background() // When using with OpenTelemetry, be sure to use the received context to propagate it.
product, err := mediator.Request[*GetProduct, *Product](ctx, query)
```

### Behaviors

Behaviors can be shared between commands and queries, but they need to be registered separately.

When registering behaviors, you must set the priority order to execute behaviors.
Priority 0 will be executed first, then 1, and so on.

```go
// Define the behavior
type LoggingBehavior struct {
  // ...
}

// Implement the IBehavior Interface
func (b *LoggingBehavior) Handle(ctx context.Context, request interface{}, next mediator.NextFunc) (interface{}, error) {
  logger.Log.Info("processing request...")
  res, err := next()
  logger.Log.Info("request processed...")
  return res, err
}

// Register the behavior for commands
behavior := &LoggingBehavior{}
order := 0
mediator.RegisterCommandBehavior(order, behavior)

// Register the behavior for queries
behavior := &LoggingBehavior{}
order := 0
mediator.RegisterQueryBehavior(order, behavior)
```

### Events

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
mediator.RegisterEventSubcriber[*ProductCreated](handler)


// Send an event synchronously
event := &ProductCreated{}
ctx := context.Background() // When using with OpenTelemetry, be sure to use the received context to propagate it.
err := mediator.PublisEvent(ctx, event)

// Send a fire and forget event
mediator.Listen() // Call this method only once in your application, like at main.go
ctx := context.Background() // When using with OpenTelemetry, be sure to use the received context to propagate it.
mediator.PublisEventAsync(ctx, event)
```
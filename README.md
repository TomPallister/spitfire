# spitfire

This is an attempt a building a somewhat generic commands, queries and events framework in go.

## How to

Example set up below..

```go
package main

import (
	"log"
	"os"
	"testing"

	"github.com/TomPallister/spitfire"
)

func main(){
    logger = log.New(os.Stdout, "Log: ", log.Ldate|log.Ltime|log.Lshortfile)
	handler = spitfire.New(logger)
}

```

## Messages
A message is a command, query or event and these can in theory be any struct but in reality they should be simple dtos.

## Commands
Each command can only have one handler. Register a command handler like below

```go
	handler.RegisterCommandHandler(createAccount{}, createAccountHandler)
```

Command handlers have the following interface

```go
    type CommandHandler = func(interface{}) (interface{}, error)
```

## Events
Each event can have multiple handlers. Register event handlers like below

```go
    handler.RegisterEventHandler(accountCreated{}, incrementAccountCreatedCount)
	handler.RegisterEventHandler(accountCreated{}, addAccountToCache)
```

Event handlers have the following interface

```go
    type EventHandler = func(interface{}) error
```

## Queries 
Each query can only have one handler. Register a query handler like below

```go
	handler.RegisterQueryHandler(getAccount{}, getAccountHandler)
```

Query handlers have the following interface

```go
    type QueryHandler = func(interface{}) (interface{}, error)
```

## Calling the handlers

In order to call your command/query and any subsequent event handlers you can do the following. Due to go's lack of
generics you need to cast the result you get back as an interface to your expected response type. Obviously how you 
call this is up to you.

```go
    result, err := handler.Handle(createAccount{UserID: 1})
	if len(err) == 0 {
		acr := result.(accountCreated)
		accountCreatedResult = &acr
	}
	errors = err
```

## Further reading
To understand spitfire fully please take a look at the test class and the code itself. It's not complex.
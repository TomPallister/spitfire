package spitfire

import (
	"errors"
	"fmt"
	"log"
)

type commandHandler struct {
	handlers map[string]func(interface{}) (interface{}, error)
	l        *log.Logger
}

func (cH *commandHandler) Register(command interface{}, handler func(interface{}) (interface{}, error)) {
	if cH.handlers == nil {
		cH.handlers = make(map[string]func(interface{}) (interface{}, error))
	}
	var key = fmt.Sprintf("%T", command)
	cH.handlers[key] = handler
}

func (cH *commandHandler) Handle(command interface{}) (interface{}, error) {
	var key = fmt.Sprintf("%T", command)
	handler, ok := cH.handlers[key]
	if !ok {
		errorMessage := fmt.Sprintf("could not find command handler for %T\n", command)
		cH.l.Printf(errorMessage)
		return nil, errors.New(errorMessage)
	}
	result, err := handler(command)
	return result, err
}

type eventHandler struct {
	handlers map[string][]func(interface{}) error
	l        *log.Logger
}

func (eH *eventHandler) Register(event interface{}, handler func(interface{}) error) {
	if eH.handlers == nil {
		eH.handlers = make(map[string][]func(interface{}) error)
	}
	var key = fmt.Sprintf("%T", event)
	//get existing handlers
	existing, ok := eH.handlers[key]
	if ok {
		new := append(existing, handler)
		eH.handlers[key] = new
	} else {
		new := make([]func(interface{}) error, 1)
		new[0] = handler
		eH.handlers[key] = new
	}
}

func (eH *eventHandler) Handle(event interface{}) (interface{}, []error) {
	var key = fmt.Sprintf("%T", event)
	handlers, ok := eH.handlers[key]
	if !ok {
		eH.l.Printf("could not find event handler for %T\n", event)
		return event, nil
	}

	errors := make([]error, 0)

	for _, h := range handlers {
		e := h(event)
		if e != nil {
			errors[len(errors)] = e
		}
	}

	return nil, errors
}

// Handler injests commands, events and queries then routes them to appropriate handlers
type Handler struct {
	eventHandler   *eventHandler
	commandHandler *commandHandler
}

// New sets up the handler and its dependencies
func New(l *log.Logger) *Handler {
	eH := &eventHandler{l: l}
	cH := &commandHandler{l: l}
	h := &Handler{eventHandler: eH, commandHandler: cH}
	return h
}

// RegisterEventHandler does what it says on the tin
func (h *Handler) RegisterEventHandler(event interface{}, handler func(interface{}) error) {
	h.eventHandler.Register(event, handler)
}

// RegisterCommandHandler does what it says on the tin
func (h *Handler) RegisterCommandHandler(command interface{}, handler func(interface{}) (interface{}, error)) {
	h.commandHandler.Register(command, handler)
}

// RegisterQueryHandler does what it says on the tin
func (h *Handler) RegisterQueryHandler(command interface{}, handler func(interface{}) (interface{}, error)) {
	h.commandHandler.Register(command, handler)
}

// Handle does what it says on the tin
func (h *Handler) Handle(command interface{}) (interface{}, []error) {
	result, err := h.commandHandler.Handle(command)
	if err != nil {
		return nil, []error{err}
	}

	_, errs := h.eventHandler.Handle(result)
	if len(errs) > 0 {
		return nil, errs
	}

	return result, nil
}
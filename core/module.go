package core

import "github.com/hjr265/aglet/cfg"

type Module interface {
	TriggerChan() <-chan struct{}
}

var ModuleFuncs = map[string]func(cfg.Hook) (Module, error){}

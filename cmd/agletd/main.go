package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hjr265/aglet/cfg"
	"github.com/hjr265/aglet/core"

	_ "github.com/hjr265/aglet/modules/xidle"
)

func main() {
	err := cfg.Load()
	if err != nil {
		log.Fatalf("load configuration: %s", err)
	}

	for k, hook := range cfg.Current.Hooks {
		mod, err := core.ModuleFuncs[hook.Type](hook)
		if err != nil {
			log.Fatalf("initialize module (%s): %s", k, err)
		}

		go func(k string, mod core.Module, hook cfg.Hook) {
			for range mod.TriggerChan() {
				log.Printf("Triggering %s", k)
			}
		}(k, mod, hook)
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}

package xidle

import (
	"log"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/hjr265/aglet/cfg"
	"github.com/hjr265/aglet/core"
)

type XIdle struct {
	duration time.Duration

	triggerChan chan struct{}
	errorChan   chan error
}

func NewXIdle(hook cfg.Hook) (core.Module, error) {
	config := struct {
		Duration cfg.Duration
	}{}
	err := hook.Decode(&config)
	if err != nil {
		return nil, err
	}

	xi := XIdle{
		duration:    time.Duration(config.Duration),
		triggerChan: make(chan struct{}),
		errorChan:   make(chan error),
	}
	xi.start()
	return &xi, nil
}

func (x XIdle) start() {
	go func() {
		conn, err := xgb.NewConn()
		if err != nil {
			x.errorChan <- err
			return
		}
		err = screensaver.Init(conn)
		if err != nil {
			x.errorChan <- err
			return
		}

		root := xproto.Setup(conn).DefaultScreen(conn).Root

		triggered := false

		for {
			reply, err := screensaver.QueryInfo(conn, xproto.Drawable(root)).Reply()
			if err != nil {
				x.errorChan <- err
				return
			}
			difference := x.duration - time.Duration(reply.MsSinceUserInput)*time.Millisecond
			elapsed := difference < 0

			if !triggered && elapsed {
				triggered = true
				x.triggerChan <- struct{}{}
			}
			if triggered && !elapsed {
				triggered = false
			}

			if !elapsed {
				log.Printf("Sleeping for %v", difference)
				time.Sleep(difference)
			} else {
				log.Printf("Sleeping for %v", -difference/2)
				time.Sleep(-difference / 2)
			}
		}
	}()
}

func (x XIdle) TriggerChan() <-chan struct{} {
	return x.triggerChan
}

func init() {
	core.ModuleFuncs["xidle"] = NewXIdle
}

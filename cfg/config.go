package cfg

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Configuation struct {
	Hooks map[string]Hook
}

var (
	Current  Configuation
	metaData toml.MetaData
)

func Load() error {
	body := struct {
		Hooks map[string]toml.Primitive
	}{}

	home := os.Getenv("HOME")
	md, err := toml.DecodeFile(filepath.Join(home, ".config", "aglet", "config.toml"), &body)
	if err != nil {
		return err
	}
	metaData = md

	Current.Hooks = map[string]Hook{}
	for k, v := range body.Hooks {
		hook := Hook{
			rest: v,
		}
		err := md.PrimitiveDecode(v, &hook)
		if err != nil {
			return err
		}
		Current.Hooks[k] = hook
	}

	return nil
}

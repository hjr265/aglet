package cfg

import "github.com/BurntSushi/toml"

type Hook struct {
	Type string
	Exec string

	rest toml.Primitive
}

func (h Hook) Decode(v interface{}) error {
	return metaData.PrimitiveDecode(h.rest, v)
}

package cfg

import "time"

type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) (err error) {
	x, err := time.ParseDuration(string(text))
	*d = Duration(x)
	return
}

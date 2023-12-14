package masksensitive

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/showa-93/go-mask"
)

type player struct {
	Name     string
	LastName string
	Phone    string `mask:"hidelast3"`
	Card     string `mask:"fixed"`
}

func init() {
	mask.SetMaskChar("*")
	var fn mask.MaskStringFunc = func(arg string, val string) (string, error) {
		var (
			n   int
			err error
		)
		if arg != "" {
			n, err = strconv.Atoi(arg)
			if err != nil {
				return "", err
			}
		}

		if n == 0 || n > len(val) {
			n = len(val)
		}

		val = val[:len(val)-n]

		val = val + strings.Repeat("*", n)

		return val, nil
	}

	mask.RegisterMaskStringFunc("hidelast", fn)
}

func (p player) Mask() player {
	t, err := mask.Mask(p)
	if err != nil {
		panic(err)
	}

	return t
}

func (p player) MarshalJSON() ([]byte, error) {
	type result struct {
		Name     string
		LastName string
		Phone    string
		Card     string
	}

	var r result

	r = result(p.Mask())

	return json.Marshal(r)
}

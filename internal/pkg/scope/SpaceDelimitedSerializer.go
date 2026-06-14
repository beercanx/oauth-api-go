package scope

import (
	"fmt"
	"strings"
)

const empty = ""
const space = " "

func marshalSpaceDelimited[T fmt.Stringer](values []T) string {

	switch len(values) {
	case 0:
		return empty
	case 1:
		return values[0].String()
	}

	var b strings.Builder
	b.WriteString(values[0].String())
	for _, s := range values[1:] {
		b.WriteString(space)
		b.WriteString(s.String())
	}
	return b.String()
}

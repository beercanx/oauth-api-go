package scope

import "strings"

const empty = ""
const space = " "

func marshalSpaceDelimited[T any](values []T, getValue func(T) string) string {

	switch len(values) {
	case 0:
		return empty
	case 1:
		return getValue(values[0])
	}

	var b strings.Builder
	b.WriteString(getValue(values[0]))
	for _, s := range values[1:] {
		b.WriteString(space)
		b.WriteString(getValue(s))
	}
	return b.String()
}

//func unmarshalScopes[T any](scope string, setValue func(string) T) []T {
//	return strings.Split(scope, space)
//}

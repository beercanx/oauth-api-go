package scope

import "strings"

const blank = ""
const space = " "

func MarshalScopes(scopes []Scope) string {

	switch len(scopes) {
	case 0:
		return blank
	case 1:
		return scopes[0].value
	}

	var b strings.Builder
	b.WriteString(scopes[0].value)
	for _, s := range scopes[1:] {
		b.WriteString(space)
		b.WriteString(s.value)
	}
	return b.String()
}

//func UnmarshalScopes(scope string) []string {
//	return strings.Split(scope, space)
//}

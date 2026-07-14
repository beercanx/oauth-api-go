package scope

func AsUniqueScopes(scopeStrings []string) Scopes {
	seen := make(map[Scope]struct{}, len(scopeStrings))
	scopes := make(Scopes, 0, len(scopeStrings))
	for _, s := range scopeStrings {
		scope := Scope(s)
		if _, ok := seen[scope]; !ok {
			seen[scope] = struct{}{}
			scopes = append(scopes, scope)
		}
	}
	return scopes
}

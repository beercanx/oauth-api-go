package scope

type Scope struct {
	value string
}

// Basic scope
//
// Deprecated: To be replaced by config sourced Clients
var Basic = Scope{"basic"}

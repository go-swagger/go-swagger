// Package defaultallofembeds is a minimal fixture exercising the --default-allof-embeds flag
// (codescan DefaultAllOfForEmbeds).
//
// Derived plainly embeds Base. By default the embed's properties are inlined into Derived; with the
// flag, Derived becomes an allOf composition ($ref to Base + its own fields).
package defaultallofembeds

// swagger:model Base
type Base struct {
	ID string `json:"id"`
}

// Derived plainly embeds Base.
//
// swagger:model Derived
type Derived struct {
	Base

	Extra string `json:"extra"`
}

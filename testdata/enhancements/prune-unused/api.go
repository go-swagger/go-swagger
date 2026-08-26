// Package pruneunused is a minimal fixture exercising the --prune flag (codescan PruneUnusedModels).
//
// Used is reachable through the route response; Orphan is a swagger:model that nothing references.
// With ScanModels and no prune, both are emitted; with prune, Orphan is dropped.
package pruneunused

import "net/http"

// Used is referenced by the route response, so it survives pruning.
//
// swagger:model Used
type Used struct {
	Name string `json:"name"`
}

// Orphan is annotated but referenced by nothing: prune drops it.
//
// swagger:model Orphan
type Orphan struct {
	Value int `json:"value"`
}

// UsedResponse carries Used.
//
// swagger:response UsedResponse
type UsedResponse struct {
	// in: body
	Body Used `json:"body"`
}

// swagger:route GET /used things GetUsed
//
// Returns a Used.
//
// Responses:
// 200: UsedResponse
func getUsed(w http.ResponseWriter, r *http.Request) {}

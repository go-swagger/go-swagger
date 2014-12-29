package swagger

// RouteParam is a object to capture route params in a framework agnostic way.
// implementations of the muxer should use these route params to communicate with the
// swagger framework
type RouteParam struct {
	Name  string
	Value string
}

// RouteParams the collection of route params
type RouteParams []RouteParam

// Get gets the value for the route param for the specified key
func (r RouteParams) Get(name string) string {
	for _, p := range r {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}

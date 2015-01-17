package spec

type resolver struct {
	root interface{}
}

// 1. check if fragment only ref, use jsonpointer to get the value
// 2. check if this ref has a url
// 3. get the data for the url
// 4. check if this data has a ref
// 5a. with a ref use this data and the new ref to do 1
// 5b. when no ref use this data to resolve the still current ref fragment

// for resolving ref uri for pointer
// check id, resolve against base
// check ref resolve against base, when ref is found exit processing

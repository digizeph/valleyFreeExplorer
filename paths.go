package main

// Depth first search for all potential paths
type Direction int

const (
	GoUp Direction = iota
	GoDown
	GoPeer
)

// depthFirstSearch recursively searches for all paths
func DepthFirstSearch(asptr *As, direction Direction, curPath Path, asesSeen map[int]bool) []Path {
	// fmt.Printf("%d\n", asptr.Asn)

	var path Path
	var paths []Path
	path = append(path, asptr)

	// loop detection
	for _, p := range curPath {
		if asptr == p {
			return []Path{}

		}
	}

	// have searched this AS previously
	if _, ok := asesSeen[asptr.Asn]; ok {
		return []Path{path}
	}

	asesSeen[asptr.Asn] = true

	switch direction {
	case GoUp:
		// this path was propagated from a customer to a provider
		// can still go up, peer, down
		for _, customer := range asptr.Customers {
			// get all paths from this customer
			newPaths := DepthFirstSearch(customer, GoDown, append(curPath, asptr), asesSeen)
			// add to paths list
			for _, ptr := range newPaths {
				paths = append(paths, append(path, ptr...))
			}
		}
		for _, customer := range asptr.Peers {
			// get all paths from this customer
			newPaths := DepthFirstSearch(customer, GoPeer, append(curPath, asptr), asesSeen)
			// add to paths list
			for _, ptr := range newPaths {
				paths = append(paths, append(path, ptr...))
			}
		}
		for _, customer := range asptr.Providers {
			// get all paths from this customer
			newPaths := DepthFirstSearch(customer, GoUp, append(curPath, asptr), asesSeen)
			// add to paths list
			for _, ptr := range newPaths {
				paths = append(paths, append(path, ptr...))
			}
		}

	default:
		// this path was propagated from a provider or a peer
		// can only go down to customers

		for _, customer := range asptr.Customers {
			// get all paths from this customer
			newPaths := DepthFirstSearch(customer, GoDown, append(curPath, asptr), asesSeen)
			// add to paths list
			for _, ptr := range newPaths {
				paths = append(paths, append(path, ptr...))
			}
		}
	}

	if len(paths) == 0 {
		paths = append(paths, path)
	}

	return paths
}


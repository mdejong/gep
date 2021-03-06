// Package functions provides a map of available functions for the GEP algorithm.
package functions

// FuncMap is a map from the symbol name of a function to its defining FuncNode.
type FuncMap map[string]FuncNode

// FuncNode defines an available function for the GEP algorithm.
type FuncNode interface {
	// Symbol is the Karva string representation of the function.
	Symbol() string
	// Terminals is the number of input terminals for the function.
	Terminals() int
	// BoolFunction represents a general boolean function.
	BoolFunction(a, b, c, d bool) bool
	// Float64Function represents a general floating-point function.
	Float64Function(a, b, c, d float64) float64
}

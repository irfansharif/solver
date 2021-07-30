package cpsatsolver

type LinearExpr = *linearExpr

type linearExpr struct {
	vars   []IntVar
	coeffs []int64
	offset int64
}

// XXX: Could instead construct a linear constraint iteratively: starting with
// bounds, setting coefficient per int var, and deciding to maximize/minimize.

func NewLinearExpr(is []IntVar, coeffs []int64, offset int64) LinearExpr {
	return &linearExpr{
		vars:   is,
		coeffs: coeffs,
		offset: offset,
	}
}

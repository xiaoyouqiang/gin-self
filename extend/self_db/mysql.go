package self_db

type SqlOp string
type M map[string]interface{}

var (
	Equal              = SqlOp("=")
	NotEqual           = SqlOp("<>")
	GreaterThan        = SqlOp(">")
	GreaterThanOrEqual = SqlOp(">=")
	SmallerThan        = SqlOp("<")
	SmallerThanOrEqual = SqlOp("<=")
	Like               = SqlOp("LIKE")
)

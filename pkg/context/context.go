package context

type (
	// CtxKey defines keys used for context values
	CtxKey string
)

// enumerates all context keys
const (
	patientKey CtxKey = "patient"
	clientKey CtxKey  = "client"
)

package context

type ctxKey string

// enumerates all context keys
const (
	PatientKey ctxKey = "patient"
	ClientKey  ctxKey = "client"
)

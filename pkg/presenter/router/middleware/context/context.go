package context

type ctxKey string

// enumerates all context keys
const (
	ClientKey  ctxKey = "client"
	PatientKey ctxKey = "patient"
	UserKey    ctxKey = "user"
)

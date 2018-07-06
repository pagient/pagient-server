package websocket

// MessageType is the type of a websocket message
type MessageType string

const (
	// MessageTypePatientUpdate marks a message that originates from a patient update operation
	MessageTypePatientUpdate MessageType = "patient_update"
	// MessageTypePatientDelete marks a message that originates from a patient delete operation
	MessageTypePatientDelete MessageType = "patient_delete"
)

// Message struct
type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

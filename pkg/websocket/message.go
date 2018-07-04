package websocket

type MessageType string

const (
	MessageTypePatientUpdate MessageType = "patient_update"
	MessageTypePatientDelete MessageType = "patient_delete"
)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

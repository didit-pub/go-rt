package rt

// StatusConstants define los estados posibles de un ticket
const (
	StatusOpen    = "open"
	StatusClosed  = "closed"
	StatusPending = "pending"
)

// ErrorResponse representa una respuesta de error de la API
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type Item struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	URL  string `json:"_url,omitempty"`
}

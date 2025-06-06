package dto

const (
	// StatusSuccess is a constant for success status
	STATUS_SUCCESS       = "success"
	MESSAGE_UNAUTHORIZED = "Unauthorized"
	// StatusError is a constant for error status
	STATUS_ERROR = "error"
)

type Response struct {
	Message    string      `json:"message"`
	Status     string      `json:"status"`
	Data       interface{} `json:"data,omitempty"`
	TotalPages int64       `json:"total_pages,omitempty"`
	// diisi oleh ResponseMeta
	Meta interface{} `json:"meta,omitempty"`
}

type ResponseMeta struct {
	AfterCursor  *string `json:"after_cursor"`
	BeforeCursor *string `json:"before_cursor"`
}

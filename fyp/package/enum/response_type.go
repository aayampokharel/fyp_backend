package enum

type RESPONSETYPE string

const (
	HTML RESPONSETYPE = "HTML"
	JSON RESPONSETYPE = "JSON"
)

func (m *RESPONSETYPE) ToString() string {
	switch *m {
	case HTML:
		return "HTML"
	case JSON:
		return "JSON"
	default:
		return ""
	}
}

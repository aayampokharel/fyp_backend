package enum

type RESPONSETYPE string

const (
	HTML     RESPONSETYPE = "HTML"
	IMAGE    RESPONSETYPE = "IMAGE"
	JSON     RESPONSETYPE = "JSON"
	PDF      RESPONSETYPE = "PDF"
	ZIP      RESPONSETYPE = "ZIP"
	PDFORZIP RESPONSETYPE = "PDFORZIP"
)

func (m *RESPONSETYPE) ToString() string {
	switch *m {
	case HTML:
		return "HTML"
	case JSON:
		return "JSON"
	case PDF:
		return "PDF"
	case ZIP:
		return "ZIP"
	case PDFORZIP:
		return "PDFORZIP"
	case IMAGE:
		return "IMAGE"
	default:
		return ""
	}
}

package enum

type HTTPMETHOD string

const (
	METHODPOST   HTTPMETHOD = "POST"
	METHODGET    HTTPMETHOD = "GET"
	METHODPUT    HTTPMETHOD = "PUT"
	METHODDELETE HTTPMETHOD = "DELETE"
)

func (m *HTTPMETHOD) ToString() string {
	switch *m {
	case METHODPOST:
		return "POST"
	case METHODGET:
		return "GET"
	case METHODPUT:
		return "PUT"
	case METHODDELETE:
		return "DELETE"
	default:
		return ""
	}
}

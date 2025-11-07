package enum

type STATUS string

const (
	PENDING  STATUS = "PENDING"
	APPROVED STATUS = "APPROVED"
	REJECTED STATUS = "REJECTED"
)

func (h STATUS) ToString() string {
	switch h {
	case PENDING:
		return "PENDING"
	case APPROVED:
		return "APPROVED"
	case REJECTED:
		return "REJECTED"
	default:
		return ""
	}
}

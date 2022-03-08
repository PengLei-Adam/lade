package contract

const IDKey = "lade:id"

type IDService interface {
	NewID() string
}

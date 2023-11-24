package events

const (
	UNKNOWN Type = iota
	MESSAGE
)

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

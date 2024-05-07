package broadcast

type Broadcast struct {
	Match
	Station Station
}

const (
	niceDate   = "Monday, Jan 2"
	timeLayout = "15:04"
)

package channel

import "whitford.io/radiofooty/internal/broadcast"

type ManualGetter struct{}

func (mg ManualGetter) GetMatches() ([]broadcast.Broadcast, error) {
	return []broadcast.Broadcast{}, nil
}

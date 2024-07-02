package application

import (
	"github.com/link-identity/app/domain"
	"sync"
)

type ILocation interface {
	GetLastNLocation(rider string, lastN int) []domain.Location
	UpdateLocation(rider string, currLocation domain.Location)
}

type location struct {
	riderLocations map[string][]domain.Location
	sync           sync.RWMutex
}

func NewLocationService() ILocation {
	return &location{}
}

func (l *location) GetLastNLocation(rider string, lastN int) []domain.Location {
	//slice := l.riderLocations[riderLocations][len()]
	//addr val addr
	var locations []domain.Location
	for i := len(l.riderLocations[rider]) - 1; i >= len(l.riderLocations[rider])-lastN && i >= 0; i-- {
		locations = append(locations, l.riderLocations[rider][i])
	}
	return locations
}

func (l *location) UpdateLocation(rider string, currLocation domain.Location) {
	if l.riderLocations == nil {
		l.riderLocations = make(map[string][]domain.Location)
	}
	if l.riderLocations[rider] == nil {
		l.riderLocations[rider] = make([]domain.Location, 0)
	}
	l.sync.Lock()
	defer func() {
		l.sync.Unlock()
	}()

	//err
	l.riderLocations[rider] = append(l.riderLocations[rider], currLocation)
}

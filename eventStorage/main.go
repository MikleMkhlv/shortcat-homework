package main

import (
	"fmt"
	"maps"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type Event struct {
	ID        int
	Type      string
	Data      string
	Timestamp time.Time
}

type EventStore struct {
	store  map[int]Event
	lastId int64
	mu     sync.Mutex
}

func NewEventStore() *EventStore {
	return &EventStore{
		store:  make(map[int]Event, 100),
		lastId: 0,
	}
}

func (es *EventStore) Add(eventType string, data string) int {
	// TODO: создать событие, добавить в хранилище, вернуть ID
	currentId := atomic.AddInt64(&es.lastId, 1)
	event := Event{
		ID:        int(currentId),
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	es.mu.Lock()
	es.store[int(currentId)] = event
	es.mu.Unlock()
	return int(currentId)
}

func (es *EventStore) GetAll() []Event {
	// TODO: вернуть копию всех событий
	es.mu.Lock()
	events := slices.Collect(maps.Values(es.store))
	es.mu.Unlock()
	return events
}

func (es *EventStore) GetByID(id int) (Event, bool) {
	// TODO: найти событие по ID
	es.mu.Lock()
	defer es.mu.Unlock()
	if value, ok := es.store[id]; ok {
		return value, true
	}
	return Event{}, false
}

func (es *EventStore) Count() int {
	es.mu.Lock()
	count := len(es.store)
	es.mu.Unlock()
	return count
}

func (es *EventStore) GetByType(eventType string) []Event {
	es.mu.Lock()
	allEvents := slices.Collect(maps.Values(es.store))
	es.mu.Unlock()

	filteredEvents := slices.DeleteFunc(allEvents, func(e Event) bool {
		return e.Type != eventType
	})
	return filteredEvents
}

func (es *EventStore) FindAfter(timestamp time.Time) []Event {
	es.mu.Lock()
	filteredEvents := make([]Event, 0, len(es.store))

	for _, e := range es.store {
		if e.Timestamp.After(timestamp) {
			filteredEvents = append(filteredEvents, e)
		}
	}
	es.mu.Unlock()

	slices.SortFunc(filteredEvents, func(a, b Event) int {
		return a.Timestamp.Compare(b.Timestamp)
	})

	return filteredEvents
}

func (es *EventStore) GetRange(startID, endID int) []Event {
	if startID > endID {
		return nil
	}

	es.mu.Lock()
	events := make([]Event, 0, len(es.store))

	for _, event := range es.store {
		if event.ID >= startID && event.ID <= endID {
			events = append(events, event)
		}
	}
	es.mu.Unlock()

	if len(events) == 0 {
		return nil
	}

	slices.SortFunc(events, func(a, b Event) int {
		return a.ID - b.ID
	})

	return events
}

func (es *EventStore) Filter(predicate func(Event) bool) []Event {
	es.mu.Lock()
	allEvents := slices.Collect(maps.Values(es.store))
	es.mu.Unlock()

	filteredEvents := make([]Event, 0, len(allEvents))
	for _, val := range allEvents {
		if predicate(val) {
			filteredEvents = append(filteredEvents, val)
		}
	}
	return filteredEvents
}

func main() {
	store := NewEventStore()

	id1 := store.Add("user.login", "user: alice")
	// id2 := store.Add("user.logout", "user: alice")

	if event, ok := store.GetByID(id1); ok {
		fmt.Printf("Event %d: %s - %s at %v\n",
			event.ID, event.Type, event.Data, event.Timestamp)
	}

	fmt.Printf("Total events: %d\n", store.Count())
}

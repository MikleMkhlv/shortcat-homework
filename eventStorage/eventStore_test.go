package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func setupEventStoregeWithEvents(t *testing.T) *EventStore {
	t.Helper()
	store := NewEventStore()

	events := []struct{ eventType, data string }{
		{"test-1", "t1"},
		{"test-2", "t2"},
		{"test-3", "t3"},
		{"test-4", "t4"},
		{"test-5", "t5"},
		{"test-6", "t6"},
		{"test-7", "t7"},
		{"test-8", "t8"},
		{"test-9", "t9"},
		{"test-10", "t10"},
	}

	for _, event := range events {
		eventId := store.Add(event.eventType, event.data)
		t.Logf("%v.setupEventStoregeWithEvents: create event. ID = %d", t.Name(), eventId)
	}

	return store
}

func TestGetAllEvent(t *testing.T) {
	eventStore := setupEventStoregeWithEvents(t)

	events := eventStore.GetAll()
	if len(events) != 10 {
		t.Errorf("%v: incorrect number of records. Expected 10 records.", t.Name())
	}

	t.Logf("%v. quantity records matches with checked value: expected value = %d. received valuse = %d", t.Name(), 10, len(events))
}

func TestGetEventByID(t *testing.T) {
	eventStore := setupEventStoregeWithEvents(t)

	testingEvents := []struct {
		id      int
		isFound bool
	}{
		{5, true},
		{1, true},
		{0, false},
		{10, true},
		{-60, false},
		{1000, false},
	}
	for _, testData := range testingEvents {
		_, exists := eventStore.GetByID(testData.id)
		if exists != testData.isFound {
			t.Errorf("%v: event with id == %d does not match test data: wait: %v, got: %v", t.Name(), testData.id, testData.isFound, exists)
		}
	}
}

func TestСorrectAssignmentOfIDForEvents(t *testing.T) {
	emptyEventStorage := NewEventStore()
	countEvents := 50

	wg := sync.WaitGroup{}

	wg.Add(countEvents)
	for i := 0; i < countEvents; i++ {
		typeEv := fmt.Sprintf("testT-%d", i)
		dataEv := fmt.Sprintf("data-from-testT-%d", i)
		go func() {
			defer wg.Done()
			emptyEventStorage.Add(typeEv, dataEv)
		}()
	}
	wg.Wait()

	if emptyEventStorage.lastId != int64(countEvents) {
		t.Errorf("%v: got lastId %d, wait lastId %d", t.Name(), emptyEventStorage.lastId, countEvents)
	}
	if len(emptyEventStorage.store) != countEvents {
		t.Errorf("%v: got len storage %d, wait len storage %d", t.Name(), len(emptyEventStorage.store), countEvents)
	}
}

func TestCorrectLengthEventInStorage(t *testing.T) {
	eventStore := setupEventStoregeWithEvents(t)
	expectedLength := 10
	receivedLength := eventStore.Count()
	if receivedLength != expectedLength {
		t.Errorf("%v: got len storage %d, wait len storage %d", t.Name(), receivedLength, expectedLength)
	}
}

func TestCorrectGetEventsByType(t *testing.T) {
	eventStorage := NewEventStore()

	expectedEvents := []struct{ eventType, data string }{
		{"race", "car-1"},
		{"race", "car-2"},
		{"race", "car-3"},
		{"race", "car-4"},
		{"race", "car-5"},
		{"race", "car-6"},
	}

	otherEvents := []struct{ eventType, data string }{
		{"rest-1", "r1"},
		{"rest-2", "r2"},
		{"rest-3", "r3"},
		{"rest-4", "r4"},
		{"rest-5", "r5"},
		{"rest-6", "r6"},
		{"rest-7", "r7"},
		{"rest-8", "r8"},
	}

	generalListEvents := make([]struct{ eventType, data string }, 0, len(expectedEvents)+len(otherEvents))
	generalListEvents = append(generalListEvents, expectedEvents...)
	generalListEvents = append(generalListEvents, otherEvents...)

	for _, event := range generalListEvents {
		eventStorage.Add(event.eventType, event.data)
	}

	expectedEventType := "race"
	expectedLengthEventListByType := len(expectedEvents)

	gotEventListByType := eventStorage.GetByType(expectedEventType)
	if len(gotEventListByType) != expectedLengthEventListByType {
		t.Errorf("%v: got len eventList by type(%s) %d, wait len eventList by type (%s) %d", t.Name(), expectedEventType, len(gotEventListByType), expectedEventType, expectedLengthEventListByType)
	}
	for _, event := range gotEventListByType {
		if event.Type != expectedEventType {
			t.Errorf("%v: got event (id=%d) with type %s, wait event (id=%d) with type %s", t.Name(), event.ID, event.Type, event.ID, expectedEventType)
		}
	}
}

func TestGetSnapshotEventsByAfterTime(t *testing.T) {
	eventStorage := NewEventStore()

	expectedEvents := []struct{ eventType, data string }{
		{"race", "car-1"},
		{"race", "car-2"},
		{"race", "car-3"},
		{"race", "car-4"},
		{"race", "car-5"},
		{"race", "car-6"},
	}
	otherEvents := []struct{ eventType, data string }{
		{"rest-1", "r1"},
		{"rest-2", "r2"},
		{"rest-3", "r3"},
		{"rest-4", "r4"},
		{"rest-5", "r5"},
		{"rest-6", "r6"},
		{"rest-7", "r7"},
		{"rest-8", "r8"},
	}

	for _, event := range otherEvents {
		eventStorage.Add(event.eventType, event.data)
	}

	expectetCurrentTime := time.Now()

	time.Sleep(time.Millisecond * 200)

	for _, event := range expectedEvents {
		eventStorage.Add(event.eventType, event.data)
	}

	gotEventListByAfterTime := eventStorage.FindAfter(expectetCurrentTime)
	if len(gotEventListByAfterTime) != len(expectedEvents) {
		t.Errorf("%v: got len eventList %d, wait len eventList %d", t.Name(), len(gotEventListByAfterTime), len(expectedEvents))
	}

	for _, event := range gotEventListByAfterTime {
		if event.Timestamp.Before(expectetCurrentTime) {
			t.Errorf("%v: event %d not from a cut", t.Name(), event.ID)
		}
	}
}

func TestGetEventListByRange(t *testing.T) {
	eventStorage := setupEventStoregeWithEvents(t)
	expectedEventyList := []struct {
		id              int
		eventType, data string
	}{
		{8, "test-8", "t8"},
		{9, "test-9", "t9"},
		{10, "test-10", "t10"},
	}

	eventList := eventStorage.GetRange(8, 10)
	if len(eventList) != len(expectedEventyList) {
		t.Errorf("%v: len eventList does not match the cut. got %d, wait %d", t.Name(), len(eventList), len(expectedEventyList))
	}
	for i, event := range eventList {
		if event.ID != expectedEventyList[i].id {
			t.Errorf("%v: the expected event (id=%d) is not in the received list", t.Name(), expectedEventyList[i].id)
		}
	}
}

func TestGetEventListByPredicateWhichFilteredById(t *testing.T) {
	eventStorage := setupEventStoregeWithEvents(t)
	expectedEventyList := []struct {
		id              int
		eventType, data string
	}{
		{8, "test-8", "t8"},
		{9, "test-9", "t9"},
		{10, "test-10", "t10"},
	}
	predicate := func(e Event) bool {
		return e.ID > 7
	}

	filteredEventList := eventStorage.Filter(predicate)
	if len(filteredEventList) != len(expectedEventyList) {
		t.Errorf("%v: len eventList does not match the cut. got %d, wait %d", t.Name(), len(filteredEventList), len(expectedEventyList))
	}
}

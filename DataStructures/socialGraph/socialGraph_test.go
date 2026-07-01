package main

import (
	"fmt"
	"testing"
)

func setupGraphWithUsers(t *testing.T, countUsers int) *Graph {
	t.Helper()
	socialGraph := NewGraph()

	testUsers := make([]struct {
		id   int
		name string
	}, 0, countUsers)

	for i := 1; i <= countUsers; i++ {
		user := struct {
			id   int
			name string
		}{
			i,
			fmt.Sprintf("testName-%d", i),
		}
		testUsers = append(testUsers, user)
	}

	for _, user := range testUsers {
		socialGraph.AddUser(user.id, user.name)
	}

	currentUserCount := socialGraph.UserCount()
	if currentUserCount != countUsers {
		t.Errorf("%v: the number of profiles does not match the stated number. got %d, wait %d", t.Name(), currentUserCount, countUsers)
	}
	return socialGraph
}

func TestAddNewUser(t *testing.T) {
	socialGraph := setupGraphWithUsers(t, 10)

	expectedUser := struct {
		id   int
		name string
	}{
		11,
		fmt.Sprintf("testName-%d", 11),
	}

	socialGraph.AddUser(expectedUser.id, expectedUser.name)
	if _, exsist := socialGraph.GetUser(expectedUser.id); !exsist {
		t.Errorf("%v: user id=%d is not exist", t.Name(), expectedUser.id)
	}

	currentUserCount := socialGraph.UserCount()
	expectedUserCount := 11
	if currentUserCount != expectedUserCount {
		t.Errorf("%v: the number of profiles does not match the stated number. got %d, wait %d", t.Name(), currentUserCount, expectedUserCount)
	}
}

func TestCreateConnectionsBetweenTwoUsers(t *testing.T) {
	socialGraph := setupGraphWithUsers(t, 10)
	expectedUsers := []struct {
		id   int
		name string
	}{
		{id: 2},
		{id: 3},
	}

	socialGraph.AddConnection(expectedUsers[0].id, expectedUsers[1].id)

	if isLinked := сheckingСonnectivity(socialGraph, expectedUsers[0].id, expectedUsers[1].id); !isLinked {
		t.Errorf("%v: there was no linked between the users. got %v, wait %v", t.Name(), isLinked, true)
	}
}

func TestRemoveConnectionsBetweenTwoUsers(t *testing.T) {
	socialGraph := setupGraphWithUsers(t, 10)
	expectedUsers := []struct {
		id   int
		name string
	}{
		{id: 2},
		{id: 3},
	}
	socialGraph.AddConnection(expectedUsers[0].id, expectedUsers[1].id)

	if isLinked := сheckingСonnectivity(socialGraph, expectedUsers[0].id, expectedUsers[1].id); !isLinked {
		t.Errorf("%v: there was no linked between the users. got %v, wait %v", t.Name(), isLinked, true)
	}

	socialGraph.RemoveConnection(expectedUsers[0].id, expectedUsers[1].id)

	if isLinked := сheckingСonnectivity(socialGraph, expectedUsers[0].id, expectedUsers[1].id); isLinked {
		t.Errorf("%v: there was linked between the users. got %v, wait %v", t.Name(), isLinked, true)
	}
}

func TestRemoveNeighborConnectionIfRemovedUser(t *testing.T) {
	socialGraph := setupGraphWithUsers(t, 10)
	expectedUsers := []struct {
		id   int
		name string
	}{
		{id: 2},
		{id: 3},
	}

	socialGraph.AddConnection(expectedUsers[0].id, expectedUsers[1].id)

	if isLinked := сheckingСonnectivity(socialGraph, expectedUsers[0].id, expectedUsers[1].id); !isLinked {
		t.Errorf("%v: there was no linked between the users. got %v, wait %v", t.Name(), isLinked, true)
	}

	socialGraph.RemoveUser(expectedUsers[0].id)

	if _, exsist := socialGraph.GetUser(expectedUsers[0].id); exsist {
		t.Errorf("%v: user id=%d is not exist", t.Name(), expectedUsers[0].id)
	}

	if isLinked := сheckingСonnectivity(socialGraph, expectedUsers[0].id, expectedUsers[1].id); isLinked {
		t.Errorf("%v: there was linked between the users. got %v, wait %v", t.Name(), isLinked, true)
	}
}

func TestCheckIsMutalWithTwousers(t *testing.T) {
	socialGraph := setupGraphWithUsers(t, 10)
	expectedUsers := []struct {
		id   int
		name string
	}{
		{id: 2},
		{id: 3},
	}

	socialGraph.AddConnection(expectedUsers[0].id, expectedUsers[1].id)
	if isLinked := socialGraph.IsMutual(expectedUsers[0].id, expectedUsers[1].id); !isLinked {
		t.Errorf("%v: there was linked between the users. got %v, wait %v", t.Name(), isLinked, true)
	}
}

func TestCreateAOneToManyRelAndCheckConnectionCount(t *testing.T) {
	socialGraph := setupGraphWithUsers(t, 10)
	expectedFrendUsers := []struct {
		id   int
		name string
	}{
		{id: 3},
		{id: 4},
		{id: 5},
		{id: 6},
	}
	expectedRootUser := struct {
		id   int
		name string
	}{
		id: 2,
	}

	for _, fUser := range expectedFrendUsers {
		socialGraph.AddConnection(expectedRootUser.id, fUser.id)
	}

	expectedConnectionCount := len(expectedFrendUsers)
	currentConnectionCount := socialGraph.ConnectionCount(expectedRootUser.id)
	if currentConnectionCount != expectedConnectionCount {
		t.Errorf("%v: the number of connections does not match. got %d, wait %d", t.Name(), currentConnectionCount, expectedConnectionCount)
	}
}

func TestGetCommonConnections(t *testing.T) {
	socialGraph := setupGraphWithUsers(t, 10)
	expectedFrendUsers := []struct {
		id   int
		name string
	}{
		{id: 2},
		{id: 3},
	}

	expectedCommonFrends := []struct {
		id   int
		name string
	}{
		{id: 5},
		{id: 6},
	}

	socialGraph.AddConnection(expectedFrendUsers[0].id, expectedFrendUsers[1].id)

	for _, f := range expectedFrendUsers {
		for _, ff := range expectedCommonFrends {
			socialGraph.AddConnection(f.id, ff.id)
		}
	}

	currentCommonConnections := socialGraph.CommonConnections(expectedFrendUsers[0].id, expectedFrendUsers[1].id)
	if len(currentCommonConnections) != len(expectedCommonFrends) {
		t.Errorf("%v: the number of common connections does not match. got %d, wait %d", t.Name(), len(currentCommonConnections), len(expectedCommonFrends))
	}
}

func TestCheckSuggestConnections(t *testing.T) {
	socialGraph := setupGraphWithUsers(t, 20)

	idFirst := 2
	idSecond := 3

	expectedFrendsOnlyFirstUser := []struct {
		id,
		frendId int
	}{
		{id: 5, frendId: 10},
		{id: 6, frendId: 11},
	}

	expectedFrendsOnlySeccondUser := []struct {
		id,
		frendId int
	}{
		{id: 8, frendId: 12},
		{id: 8, frendId: 13},
		{id: 7, frendId: 14},
		{id: 7, frendId: 15},
	}

	socialGraph.AddConnection(idFirst, idSecond)

	// Настраиваем связи первого уровня
	createConnectionsOneToMany(socialGraph, idFirst, expectedFrendsOnlyFirstUser)
	createConnectionsOneToMany(socialGraph, idSecond, expectedFrendsOnlySeccondUser)

	// Настраиваем связи второго уровня (друзья друзей)
	for _, conn := range expectedFrendsOnlyFirstUser {
		socialGraph.AddConnection(conn.id, conn.frendId)
	}
	for _, conn := range expectedFrendsOnlySeccondUser {
		socialGraph.AddConnection(conn.id, conn.frendId)
	}

	friendsOfFirst := socialGraph.SuggestConnections(idFirst)
	if len(friendsOfFirst) != 4 {
		t.Errorf("for ID %d 4 recommendations were expected, received %d", idFirst, len(friendsOfFirst))
	}

	friendsOfSecond := socialGraph.SuggestConnections(idSecond)
	if len(friendsOfSecond) != 6 {
		t.Errorf("for ID %d 6 recommendations were expected, received %d", idSecond, len(friendsOfSecond))
	}
}

func сheckingСonnectivity(socialGraph *Graph, fromID, toID int) bool {
	return socialGraph.HasConnection(fromID, toID) && socialGraph.HasConnection(toID, fromID)
}

func createConnectionsOneToMany(socialGraph *Graph, fromId int, friend []struct {
	id      int
	frendId int
}) {
	for _, fUser := range friend {
		socialGraph.AddConnection(fromId, fUser.id)
	}
}

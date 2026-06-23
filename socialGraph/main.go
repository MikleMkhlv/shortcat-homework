package main

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type User struct {
	ID    int
	Name  string
	Links map[int]Connection
}

type Graph struct {
	data  map[int]map[int]*User
	users map[int]*User
	mx    sync.Mutex
}

func NewGraph() *Graph {
	return &Graph{
		data:  make(map[int]map[int]*User, 100),
		users: make(map[int]*User),
	}
}

func (g *Graph) AddUser(id int, name string) {
	g.mx.Lock()
	defer g.mx.Unlock()
	if _, ok := g.data[id]; !ok {
		g.data[id] = make(map[int]*User, 10)
		g.users[id] = &User{ID: id, Name: name, Links: make(map[int]Connection)}
	}
}

func (g *Graph) GetUser(id int) (*User, bool) {
	g.mx.Lock()
	defer g.mx.Unlock()
	user, ok := g.users[id]
	return user, ok
}

func (g *Graph) AddConnection(fromID, toID int) bool {
	g.mx.Lock()
	defer g.mx.Unlock()
	if fromID == toID {
		return false
	}

	user1, ok1 := g.users[fromID]
	user2, ok2 := g.users[toID]

	if !ok1 || !ok2 {
		return false
	}

	if _, exists := g.data[fromID][toID]; exists {
		return false
	}

	g.data[user1.ID][user2.ID] = user2
	g.data[user2.ID][user1.ID] = user1

	return true
}

func (g *Graph) GetConnections(userID int) []*User {
	return slices.Collect(maps.Values(g.data[userID]))
}

func (g *Graph) HasConnection(fromID, toID int) bool {
	if _, ok := g.data[fromID][toID]; ok {
		return true
	}
	return false
}

func (g *Graph) UserCount() int {
	g.mx.Lock()
	defer g.mx.Unlock()
	userCount := len(g.users)
	return userCount
}

func (g *Graph) RemoveConnection(fromID, toID int) bool {
	g.mx.Lock()
	defer g.mx.Unlock()
	if fromID == toID {
		return false
	}

	user1, ok1 := g.users[fromID]
	user2, ok2 := g.users[toID]

	if !ok1 || !ok2 {
		return false
	}

	delete(g.data[user1.ID], user2.ID)
	delete(g.data[user2.ID], user1.ID)

	return true
}

func (g *Graph) RemoveUser(id int) bool {
	g.mx.Lock()
	defer g.mx.Unlock()
	if _, ok := g.users[id]; !ok {
		return false
	}
	for _, val := range g.data[id] {
		if g.HasConnection(val.ID, id) {
			delete(g.data[val.ID], id)
		}
	}
	delete(g.data, id)
	delete(g.users, id)

	return true
}

func (g *Graph) IsMutual(id1, id2 int) bool {
	g.mx.Lock()
	defer g.mx.Unlock()
	_, ok1 := g.data[id1][id2]
	_, ok2 := g.data[id2][id1]
	return ok1 || ok2
}

func (g *Graph) ConnectionCount(userID int) int {
	return len(g.data[userID])
}

func (g *Graph) CommonConnections(id1, id2 int) []*User {
	g.mx.Lock()
	defer g.mx.Unlock()
	// TODO: найти пользователей, с которыми связаны оба
	general := make([]*User, 0, len(g.users))
	if id1 == id2 {
		return nil
	}
	for key, user := range g.data[id1] {
		if key == id2 {
			continue
		}

		if _, ok := g.data[id1][id2]; !ok {
			return nil
		}
		// if g.HasConnection(id2, key) {
		// 	general = append(general, user)
		// }
		general = append(general, user)

	}
	return general
}

func (g *Graph) SuggestConnections(userID int) []*User {
	// TODO: найти друзей друзей, исключая текущие связи и самого пользователя
	// g.mx.Lock()
	// defer g.mx.Unlock()
	// frendList := make([]*User, 0, 30)
	// for fid := range g.data[userID] {
	// 	for ffid, ffUser := range g.data[fid] {
	// 		if g.HasConnection(userID, ffid) || g.HasConnection(fid, ffid) {
	// 			continue
	// 		}
	// 		frendList = append(frendList, ffUser)
	// 	}
	// }
	// return frendList
	g.mx.Lock()
	defer g.mx.Unlock()
	userFrends, exists := g.data[userID]
	if !exists {
		return nil
	}

	excluded := make(map[int]bool, len(userFrends))
	excluded[userID] = true
	for fid := range userFrends {
		excluded[fid] = true
	}

	verifiedSetFrend := make(map[int]bool)
	friendsList := make([]*User, 0, 30)

	for fid := range userFrends {
		for ffid, user := range g.data[fid] {
			if excluded[ffid] || verifiedSetFrend[ffid] {
				continue
			}
			verifiedSetFrend[ffid] = true
			friendsList = append(friendsList, user)
		}
	}

	return friendsList

}

func (g *Graph) GetAllUsers() []*User {
	return slices.Collect(maps.Values(g.users))
}

func main() {
	graph := NewGraph()

	graph.AddUser(1, "Alice")
	graph.AddUser(2, "Bob")
	graph.AddUser(3, "Charlie")

	graph.AddConnection(1, 2) // Alice -> Bob
	graph.AddConnection(1, 3) // Alice -> Charlie
	graph.AddConnection(2, 3) // Bob -> Charlie

	// fmt.Printf("delete user Alice: %v\n", graph.RemoveUser(1))
	fmt.Printf("ismutal user Alice and Charlie: %v\n", graph.IsMutual(1, 3))
	graph.AddTypedConnection(1, 2, Friend{})

	typeF, ok := graph.GetConnectionInfo(1, 2)
	if ok {
		fmt.Printf("connection Alice with Bob: %v - %s\n", ok, typeF.Type())
	} else {

	}
	if user, ok := graph.GetUser(1); ok {
		fmt.Printf("User: %s\n", user.Name)
		friends := graph.GetConnections(1)
		fmt.Printf("Friends: %d\n", len(friends))
		for _, friend := range friends {
			fmt.Printf("  - %s\n", friend.Name)
		}
	}

	fmt.Printf("Alice and Bob connected: %v\n",
		graph.HasConnection(1, 2))

}

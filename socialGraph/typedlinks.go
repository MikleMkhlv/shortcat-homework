package main

type Connection interface {
	Type() string
	Weight() int
}

type Friend struct {
	Since string // Дата начала дружбы
}

func (f Friend) Type() string {
	return "friend"
}

func (f Friend) Weight() int {
	return 10
}

type Follower struct {
	Notifications bool
}

func (f Follower) Type() string {
	return "follower"
}

func (f Follower) Weight() int {
	return 5
}

type Blocked struct {
	Reason string
}

func (b Blocked) Type() string {
	return "blocked"
}

func (b Blocked) Weight() int {
	return -1
}

func (g *Graph) AddTypedConnection(fromID, toID int, conn Connection) bool {
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

	user1.Links[toID] = conn
	user2.Links[fromID] = conn

	g.data[user1.ID][user2.ID] = user2
	g.data[user2.ID][user1.ID] = user1

	return true
}

func (g *Graph) GetConnectionsByType(userID int, connType string) []*User {
	g.mx.Lock()
	defer g.mx.Unlock()
	// user, ok := g.users[userID]
	// if !ok {
	// 	return nil
	// }
	userFrends, exists := g.data[userID]
	if !exists {
		return nil
	}
	friendListWithTheRightType := make([]*User, 0, len(userFrends))
	for _, frend := range userFrends {
		if frend.Links[userID].Type() == connType {
			friendListWithTheRightType = append(friendListWithTheRightType, frend)
		}
	}

	return friendListWithTheRightType
}

func (g *Graph) GetConnectionInfo(fromID, toID int) (Connection, bool) {
	// TODO
	g.mx.Lock()
	defer g.mx.Unlock()

	if fromID == toID {
		return nil, false
	}

	user1, ok1 := g.users[fromID]
	user2, ok2 := g.users[toID]

	if !ok1 || !ok2 {
		return nil, false
	}

	if _, exists := g.data[fromID][toID]; !exists {
		return nil, false
	}

	if user1.Links[toID].Type() != user2.Links[fromID].Type() {
		return nil, false
	}

	return user1.Links[toID], true
}

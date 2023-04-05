package model

import "github.com/google/uuid"

var (
	AvailableLobbies []*Lobby
	StartedLobbies   []*Lobby
)

// ==== LOBBY FEATURES ======
// Create
// Delete
// AddUser
//

// Lobby defines a group of maximum 5 users
type Lobby struct {
	ID     uuid.UUID `json:"id"`
	Leader *User     `json:"leader_id"`
	Users  []*User   `json:"users"`
}

func (l *Lobby) AddUser(u *User) {
	l.Users = append(l.Users, u)
}

// Handler verify if user is the lobby leader
func (l *Lobby) RemoveUser(u *User) {
	for i, v := range l.Users {
		if v == u {
			removeUser(l.Users, i)
		}
	}
}

// Constructor function for create Lobby
func newLobby() *Lobby {
	return &Lobby{}
}

func (l *Lobby) SetLeader(user *User) {
	l.Leader = user
}

func (u *User) CreateLobby() *Lobby {
	l := newLobby()
	l.SetLeader(u)
	l.Users = append(l.Users, u)

	AvailableLobbies = append(AvailableLobbies, l)

	return l
}

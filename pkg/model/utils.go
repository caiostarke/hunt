package model

func removeQueue(slice []*Queue, s int) []*Queue {
	return append(slice[:s], slice[s+1:]...)
}

func removeUser(slice []*User, s int) []*User {
	return append(slice[:s], slice[s+1:]...)
}

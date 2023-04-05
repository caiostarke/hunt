package model

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	AvailableQueues []*Queue
	QueuesStarted   []*Queue
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type Queue struct {
	ID        uuid.UUID
	ConnUsers []User
}

func NewQueue() *Queue {
	return &Queue{
		ID: uuid.New(),
	}
}

func (q *Queue) AddUser(u User) error {
	if len(q.ConnUsers) == 9 {
		return fmt.Errorf("Queue limit exceeded")
	}

	q.ConnUsers = append(q.ConnUsers, u)
	return nil
}

func SearchQueue(u *User, logger *zap.Logger) uuid.UUID {
	if len(AvailableQueues) == 0 {
		q := NewQueue()
		q.AddUser(*u)
		AvailableQueues = append(AvailableQueues, q)
		fmt.Println(White + " ------------------------------------------------------------------------------ " + Reset)
		logger.Info(White + q.ID.String() + "Queue created " + u.ID.String() + " Added into queue" + Reset)
		fmt.Println(White + " ------------------------------------------------------------------------------ " + Reset)
		return q.ID
	} else {
		cleaner()
	}

	for _, q := range AvailableQueues {
		q.ConnUsers = append(q.ConnUsers, *u)
		logger.Warn(Blue + " " + strconv.Itoa(len(q.ConnUsers)) + " users in current queue" + Reset)
		fmt.Println(White + " ------------------------------------------------------------------------------ " + Reset)
		logger.Info(White + " " + q.ID.String() + " Queue " + u.ID.String() + " Added into queue " + Reset)
		fmt.Println(White + " ------------------------------------------------------------------------------ " + Reset)
		return q.ID
	}

	return uuid.UUID{}
}

func cleaner() {
	for i, q := range AvailableQueues {
		if len(q.ConnUsers) == 10 {
			AvailableQueues = removeQueue(AvailableQueues, i)
			QueuesStarted = append(QueuesStarted, q)
		}
	}

	if len(AvailableQueues) == 0 {
		q := NewQueue()
		AvailableQueues = append(AvailableQueues, q)
	}
}

// returns the total of users in available queues
func GetUsersFromAvailableQueues() int {
	users := 0
	for _, q := range AvailableQueues {
		users += len(q.ConnUsers)
	}

	return users
}

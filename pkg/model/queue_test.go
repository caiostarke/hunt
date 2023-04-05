package model

import (
	"testing"

	"go.uber.org/zap"
)

func TestSearchQueue(t *testing.T) {
	logger, _ := zap.NewProduction()

	users := []*User{}

	// Create 15 users and put inside []users
	for i := 0; i < 15; i++ {
		u := NewUser()
		users = append(users, u)
	}

	for _, u := range users {
		SearchQueue(u, logger)
	}

	t.Run("Testing overrun issues", func(t *testing.T) {
		if len(AvailableQueues[0].ConnUsers) > 9 {
			t.Errorf("Expecting 5 users got %v", len(AvailableQueues[0].ConnUsers))
		}
	})

	t.Run("Expect 5 users in queue", func(t *testing.T) {
		if len(AvailableQueues[0].ConnUsers) != 5 {
			t.Errorf("Expecting 5 users got %v", len(AvailableQueues[0].ConnUsers))
		}
	})

	t.Run("Check if after exceed users limit in queue, queue is dropped from AvailableQueues", func(t *testing.T) {
		if len(AvailableQueues) != 1 {
			t.Errorf("Expecting 1 queue, got %v", len(AvailableQueues))
		}
	})

}

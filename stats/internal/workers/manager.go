package workers

import (
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Start() {
	go StartViewsListener()
	go StartLikesListener()
	go StartCommentsListener()
}
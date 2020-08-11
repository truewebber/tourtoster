package repository

import (
	"sync"

	"github.com/truewebber/tourtoster/token"
)

type (
	Memory struct {
		data map[string]*token.Token
		m    sync.Mutex
	}
)

func NewMemory() *Memory {
	return &Memory{
		data: make(map[string]*token.Token),
	}
}

func (m *Memory) Token(token string) (*token.Token, error) {
	m.m.Lock()
	defer m.m.Unlock()

	if t, ok := m.data[token]; ok {
		return t, nil
	}

	return nil, nil
}

func (m *Memory) Save(token *token.Token) error {
	m.m.Lock()
	defer m.m.Unlock()

	m.data[token.Token] = token

	return nil
}

func (m *Memory) Delete(token string) error {
	m.m.Lock()
	defer m.m.Unlock()

	delete(m.data, token)

	return nil
}

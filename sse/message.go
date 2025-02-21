package sse

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Message - SSE message
type Message struct {
	event   string
	records map[string]string
	mutex   sync.Mutex
}

// AddRecord - adds new record to the message
func (m *Message) AddRecord(key string, value string) *Message {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.records[key] = value
	return m
}

// SetEvent - sets event type
func (m *Message) SetEvent(event string) *Message {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.event = event
	return m
}

// RemoveRecord - removes a record from the message
func (m *Message) RemoveRecord(key string) *Message {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.records, key)
	return m
}

// Text - converts the Message into text
func (m *Message) Bytes() []byte {
	var bytes []byte
	if len(m.event) > 0 {
		bytes = append(bytes, fmt.Sprintf("event: %s\n", m.event)...)
	}
	s, _ := json.Marshal(m.records)
	bytes = append(bytes, fmt.Sprintf("data: %s\n\n", s)...)
	return bytes
}

func NewMessage() *Message {
	return &Message{"", make(map[string]string), sync.Mutex{}}
}

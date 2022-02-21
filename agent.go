package agent

import (
	"errors"
	"log"

	"github.com/beldmian/light"
)

type MessageHandler light.EventHandler
type Event = light.Event

func GetMessage(event Event) Message {
	return event.Payload["message"].(Message)
}

type Agent interface {
	Start() error
	GetAddress() string
	MessageHandler(Event) error
}

type MessageType int

const (
	Request MessageType = iota
	Response
	CTF
	Propose
	AcceptProposal
)

type Message struct {
	FromAddr string
	ToAddr   string
	Message  string
	Type     MessageType
	Data     map[string]interface{}
}

type Manager struct {
	registeredAgents map[string]Agent
	disposer         light.Disposer
}

func NewManager() Manager {
	return Manager{
		registeredAgents: make(map[string]Agent),
		disposer:         light.NewDisposer(),
	}
}

func (m *Manager) RegisterAgent(agent Agent) {
	if _, ok := m.registeredAgents[agent.GetAddress()]; ok {
		return
	}
	m.disposer.Handle(agent.GetAddress(), agent.MessageHandler)
	m.registeredAgents[agent.GetAddress()] = agent
}

func (m *Manager) SendMessage(message Message) error {
	if _, ok := m.registeredAgents[message.ToAddr]; !ok {
		return errors.New("not found agent with address provided in ToAddr")
	}
	log.Println("Sending message to", message.ToAddr)
	return m.disposer.Emit(Event{
		Name: message.ToAddr,
		Payload: map[string]interface{}{
			"message": message,
		},
	})
}

func (m *Manager) Start() {
	endedAgents := len(m.registeredAgents)
	isEnded := make(chan int, endedAgents)
	for _, agent := range m.registeredAgents {
		go func(agent Agent) {
			if err := agent.Start(); err != nil {
				panic(err)
			}
			isEnded <- 0
		}(agent)
	}
	for range isEnded {
		endedAgents -= 1
		if endedAgents == 0 {
			break
		}
	}
}

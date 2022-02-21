package agent_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/beldmian/agent"
)

type AgentA struct {
	Addr string
}

func (a AgentA) Start() error {
	fmt.Println("Hello world")
	return nil
}

func (a AgentA) MessageHandler(event agent.Event) error {
	fmt.Println(agent.GetMessage(event).Message)
	return nil
}

func (a AgentA) GetAddress() string {
	return a.Addr
}

type AgentB struct {
	Manager agent.Manager
	Addr    string
}

func (a AgentB) Start() error {
	time.Sleep(time.Second)
	log.Println("sending")
	a.Manager.SendMessage(agent.Message{
		ToAddr:  "a",
		Message: "Hello, world",
	})
	return nil
}

func (a AgentB) MessageHandler(event agent.Event) error {
	fmt.Println(agent.GetMessage(event).Message)
	return nil
}

func (a AgentB) GetAddress() string {
	return a.Addr
}

func TestManager_Start(t *testing.T) {
	manager := agent.NewManager()
	startAgent := AgentA{
		Addr: "a",
	}
	startAgent2 := AgentA{
		Addr: "b",
	}
	manager.RegisterAgent(startAgent)
	manager.RegisterAgent(startAgent2)
	manager.Start()
}

func TestManager_SendMessage(t *testing.T) {
	manager := agent.NewManager()
	startAgent := AgentA{
		Addr: "a",
	}
	startAgent2 := AgentB{
		Addr:    "b",
		Manager: manager,
	}
	manager.RegisterAgent(startAgent)
	manager.RegisterAgent(startAgent2)
	manager.Start()
}

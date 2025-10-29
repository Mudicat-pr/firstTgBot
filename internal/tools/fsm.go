package tools

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateFunc func(msg *tgbotapi.Message) error

// Структура для работы с состояниями
type FSM struct {
	// Мапа состояний для работы с ботом
	states map[string]StateFunc

	// Юзер и его текущее состояние
	currentState map[int64]string
}

func New() *FSM {
	return &FSM{
		states:       make(map[string]StateFunc),
		currentState: make(map[int64]string),
	}
}

func BindState[T any](recv T, method func(T, *tgbotapi.Message) error) StateFunc {
	return func(msg *tgbotapi.Message) error {
		return method(recv, msg)
	}
}

func (f *FSM) Register(state string, handle StateFunc) *FSM {
	f.states[state] = handle
	return f
}

func (f *FSM) SetState(userID int64, state string) {
	f.currentState[userID] = state
}

func (f *FSM) State(userID int64) string {
	return f.currentState[userID]
}

func (f *FSM) ClearState(userID int64) {
	delete(f.currentState, userID)
}

func (f *FSM) Handle(msg *tgbotapi.Message) {
	userID := msg.From.ID
	state := f.State(userID)
	//if state == "" {
	//	state = "" // Временно
	//	f.SetState(userID, state)
	//}

	if handler, ok := f.states[state]; ok {
		err := handler(msg)
		if err != nil {
			log.Println("Failed to use handler. METHOD handle from tools/fsm.go")
			return
		}
	}
}

package tools

import (
	"errors"
	"log"

	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateFunc func(msg *tgbotapi.Message, data interface{}) (
	NextState int, UserData interface{}, err error)

type FSM struct {
	currentState map[int64]int         // Текущий статус. Крепится за юзером
	data         map[int64]interface{} // Данные. Аналогично статусу
	handle       map[int]StateFunc     // Хендлер для обработки шагов
}

func New() *FSM {
	return &FSM{
		currentState: make(map[int64]int),
		data:         make(map[int64]interface{}),
		handle:       make(map[int]StateFunc),
	}
}

func (f *FSM) SetState(userID int64, state int) {
	log.Printf("=== SetState: user=%d, state=%d ===", userID, state)
	f.currentState[userID] = state
}

func (f *FSM) UserState(userID int64) int {
	state := f.currentState[userID]
	log.Printf("=== UserState: user=%d, returns=%d ===", userID, state)
	return f.currentState[userID]
}

func (f *FSM) Register(state int, handle StateFunc) *FSM {
	f.handle[state] = handle
	log.Printf("=== Register: state=%d, handler=%p ===", state, handle)
	return f
}

func (f *FSM) ClearState(userID int64) {
	delete(f.currentState, userID)
	delete(f.data, userID)
}

func BindState[T any](recv T,
	method func(T,
		*tgbotapi.Message,
		interface{}) (int, interface{}, error)) StateFunc {
	return func(msg *tgbotapi.Message, data interface{}) (int, interface{}, error) {
		return method(recv, msg, data)
	}
}

func (f *FSM) HandleState(msg *tgbotapi.Message) (err error) {

	defer func() { err = e.WrapIfErr("Check HandleState in FSM", err) }()
	userID := msg.From.ID
	state := f.UserState(userID)
	log.Printf("=== HandleState: looking for handler for state %d ===", state)
	handle, ok := f.handle[state]
	log.Printf("=== Handler for state %d: found=%v, handler=%p ===", state, ok, handle)
	if !ok {
		log.Printf("Unknown state for handle %v", handle)
		log.Printf("=== NO HANDLER for state %d ===", state)
		return errors.New("Can't found handler for this state")
	}
	currentData := f.data[userID]
	nextState, newData, err := handle(msg, currentData)
	log.Printf("Handler returned: nextState=%d, err=%v", nextState, err)
	if err != nil {
		return err
	}
	f.data[userID] = newData

	if nextState == 0 {
		f.ClearState(userID)
	} else {
		f.SetState(userID, nextState)
	}
	return nil
}

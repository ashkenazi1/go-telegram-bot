package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type StateHandler func(bot *Bot, message *tgbotapi.Message, userState *UserState)

type UserState struct {
	State string
}

type StateManager struct {
	userStates    map[int64]*UserState
	stateHandlers map[string]StateHandler
}

func NewStateManager() *StateManager {
	return &StateManager{
		userStates:    make(map[int64]*UserState),
		stateHandlers: make(map[string]StateHandler),
	}
}

func (sm *StateManager) GetState(userID int64) *UserState {
	state, exists := sm.userStates[userID]
	if !exists {
		state = &UserState{State: "start"} // Default start state
		sm.userStates[userID] = state
	}
	return state
}

func (sm *StateManager) SetState(userID int64, state string) {
	sm.userStates[userID].State = state
}

func (sm *StateManager) RegisterState(state string, handler StateHandler) {
	sm.stateHandlers[state] = handler
}

func (sm *StateManager) GetHandler(state string) (StateHandler, bool) {
	handler, exists := sm.stateHandlers[state]
	return handler, exists
}

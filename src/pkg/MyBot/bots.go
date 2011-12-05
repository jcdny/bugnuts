package MyBot

import (
	. "bugnuts/state"
	. "bugnuts/parameters"
)

type ABot struct {
	Key    string
	Desc   string
	PKey   string
	NewBot func(*State, *Parameters) Bot
}

var botList = make(map[string]ABot)

func RegisterABot(bot ABot) {
	botList[bot.Key] = bot
}

func UnregisterABot(key string) {
	botList[key] = ABot{}, false
}

func BotList() []string {
	bl := make([]string, 0, len(botList))
	for k, b := range botList {
		bl = append(bl, k+": "+b.Desc)
	}

	return bl
}

func BotGet(k string) *ABot {
	b, ok := botList[k]
	if !ok {
		return nil
	}

	return &b
}

func NewBot(k string, s *State) Bot {
	b, ok := botList[k]
	if !ok {
		return nil
	}
	nbot := b.NewBot(s, ParameterSets[b.PKey])

	return nbot
}

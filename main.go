package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"godrunk/config"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Bot struct {
	*tgbotapi.BotAPI
}

type Card struct {
	Name        string
	Description string
}

var cards = []Card{
	{"Я", "Этот бокал для тебя."},
	{"Ты", "Выбери, кто будет пить."},
	{"Джентльмены", "Все джентльмены за столом пьют."},
	{"Леди", "Все леди за столом пьют."},
	{"Тост", "Произносишь тост, все пьют."},
	{"Передай налево", "Игрок слева от тебя пьет."},
	{"Передай направо", "Игрок справа от тебя пьет."},
	{"Вызов", "Выбери игрока и выпей. Он должен выпить не меньше тебя."},
	{"Все на пол", "Все игроки должны коснуться пола. Последний пьет."},
	{"Прозвище", "Придумай прозвище игроку. Все должны звать его этим прозвищем до конца игры. Тот, кто обратился к нему иначе, пьет."},
	{"Твои правила", "Становишься Rulemaster'ом и придумываешь свои правила. Следующий Rulemaster может отменить правила предыдущего."},
	{"Секретная служба", "Все должны приложить ладонь к уху, изображая телохранителей. Последний становится президентом и пьет."},
	{"Брудершафт", "Выпей на брудершафт с игроком на выбор."},
	{"Шах и мат", "Выбери игрока, который будет пить с тобой каждый раз, когда ты ошибаешься."},
	{"Повтори за мной", "Просишь игрока повторить за тобой (скороговорку или сложнопроизносимое слово). Если у него не получилось - он пьет. Получилось - пьешь ты."},
	{"Неудобные вопросы", "Каждый игрок имеет право задать тебе любой вопрос. Если ты отказываешься на него отвечать - ты пьешь."},
	{"Нос", "Все игроки должны коснуться носа. Последний пьет."},
	{"Категория", "Вытянувший карту придумывает категорию (марки презервативов, музыкальные группы 90-х годов, модели Mercedes). Остальные игроки называют слова из этой категории. Кто не сможет - пьет."},
	{"Я никогда не", "Говоришь то, что ты \"Никогда не делал\" (но на самом деле делал или очень хотел бы). Тот, кто делал это, пьет."},
	{"Вопросы", "Игрок задает вопрос игроку слева. Отвечать на него нельзя, нужно быстро задать вопрос следующему соседу. Сбился? Ошибся? Запнулся? Выпей."},
	{"Цвет", "Игрок называет цвет, следующий повторяет его и добавляет свой, и так далее. Кто сбился, тот пьет."},
	{"Кубок", "Первые три игрока, вытянувшие эту карту, сливают содержимое своих бокалов в кубок. Четвертый это дело выпивает."},
	{"Саймон говорит", "Тот, кто вытянул эту карту делает какой-нибудь жест, следующий делает то же самое и добавляет свой. Так продолжается, пока кто-нибудь не собьется."},
	{"Товарищ заебал", "Игрок, вытянувший эту карту, становится товарищем. Другим игрокам нельзя отвечать на его вопросы."},
}

var decks = make(map[int64][]*Card)

func createDeck() []*Card {
	maxCards := 5
	res := make([]*Card, 0, len(cards) * maxCards)
	for i := 0; i < len(cards); i++ {
		for j := 0; j < maxCards; j++ {
			res = append(res, &cards[i])
		}
	}

	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})

	return res
}

func (z *Bot) handleUpdate(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	deck, exists := decks[chatID]
	if !exists || len(deck) == 0 {
		deck = createDeck()
		decks[chatID] = deck
	}

	card := deck[0]
	decks[chatID] = deck[1:]

	text := fmt.Sprintf("*%v*\n%v", card.Name, card.Description)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	_, err := z.Send(msg)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	cfg := config.GetConfig("godrunk.yaml")

	b, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}
	bot := &Bot{b}

	bot.Debug = cfg.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	webhookConfig := tgbotapi.NewWebhook(cfg.WebhookAddress)
	_, err = bot.SetWebhook(webhookConfig)
	if err != nil {
		log.Panic(err)
	}

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil)

	for update := range updates {
		bot.handleUpdate(update)
	}
}

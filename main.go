package main

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"math/rand"
	"strings"
	"time"
)

func main() {

	// Хендлер бота (vk) и лонгполла (lp)
	vk := api.NewVK(TOKEN)
	group, _ := vk.GroupsGetByID(nil)
	lp, _ := longpoll.NewLongPoll(vk, group[0].ID)

	rand.Seed(time.Now().UnixNano())

	// Событие Новое сообщение
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {

		// Получение POSIX-времени из GetPOSIX.go
		var posix = getPOSIX()

		// Логируем сообщение
		log.Printf("%d %d: %s", obj.Message.PeerID, obj.Message.FromID, obj.Message.Text)

		// Перевод сообщение в нижний регистр для последующего поиска в нем
		obj.Message.Text = strings.ToLower(obj.Message.Text)

		if obj.Message.Text == "расписос" {
			// Собираем сообщение-ответ
			b := params.NewMessagesSendBuilder()
			b.Message(getSchedule(posix, obj.Message.FromID, obj.Message.PeerID))
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			vk.MessagesSend(b.Params)
		}

		if strings.Contains(obj.Message.Text, "расписос на завтра") {
			b := params.NewMessagesSendBuilder()
			b.Message(getSchedule(posix+86400, obj.Message.FromID, obj.Message.PeerID))
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			vk.MessagesSend(b.Params)
		}

		if strings.Contains(obj.Message.Text, "тупой бот") || strings.Contains(obj.Message.Text, "бот тупой") {
			b := params.NewMessagesSendBuilder()
			b.Attachment("photo-208113987_457239017")
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			vk.MessagesSend(b.Params)
		}

		if strings.Contains(obj.Message.Text, "анекдот от марченко") {
			jokeNumber := rand.Intn(4-0) + 0
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message(jokes[jokeNumber])
			b.PeerID(obj.Message.PeerID)
			vk.MessagesSend(b.Params)
		}

		if obj.Message.Text == "ты где" {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("братка я в доту хуярю")
			b.PeerID(obj.Message.PeerID)
			vk.MessagesSend(b.Params)
		}

		if obj.Message.Text == "posix" {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message(fmt.Sprint(posix))
			b.PeerID(obj.Message.PeerID)
			vk.MessagesSend(b.Params)
		}
	})

	// Запуск lp-хендлера
	err := lp.Run()
	if err != nil {
		log.Fatal(err)
	}
}

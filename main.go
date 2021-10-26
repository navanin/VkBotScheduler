package main

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"strings"
)

func main() {

	// Хендлер бота (vk) и лонгполла (lp)
	vk := api.NewVK(TOKEN)
	group, _ := vk.GroupsGetByID(nil)
	lp, _ := longpoll.NewLongPoll(vk, group[0].ID)

	// Событие Новое сообщение
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {

		// Получение POSIX-времени из GetPOSIX.go
		var posix = getPOSIX()

		// Логируем сообщение
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		// Перевод сообщение в нижний регистр для последующего поиска в нем
		obj.Message.Text = strings.ToLower(obj.Message.Text)

		if strings.Contains(obj.Message.Text, "расписос") {
			// Собираем сообщение-ответ
			b := params.NewMessagesSendBuilder()
			b.Message(getSchedule(posix))
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			vk.MessagesSend(b.Params)
		}

		if strings.Contains(obj.Message.Text, "расписос на завтра") {
			b := params.NewMessagesSendBuilder()
			b.Message(getSchedule(posix + 86400))
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

	})

	// Запуск lp-хендлера
	err := lp.Run()
	if err != nil {
		log.Fatal(err)
	}
}

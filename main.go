package main

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"time"
)

var posix = 1634947200 								// Время в формате POSIX, ибо другое дневник, пидорас эдакий, не жрет. 23 октября 2021.
var latestDay = 23 									// Последний день, по мнению бота - 23 октября 2021.

func main() {

	vk := api.NewVK(TOKEN) 							// указание токена для работы бота и создание хендлера
	group, _ := vk.GroupsGetByID(nil) 		// Указание группы, от которой пойдут сообщения
	lp, _ := longpoll.NewLongPoll(vk, group[0].ID)	// Создание Long Poll'a

	// Событие "Новое сообщение"
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text) //логируем сообщение

		todayDay := int(time.Now().Day()) 	// Запись сегодняшнего числа в переменную
		if latestDay != todayDay {  		// Проверка, совпадает ли сегодняшний день с последним днем работы бота
			var diff = latestDay - todayDay // получаем разницу в днях, на случай того, если ботом не пользовались больше 2-х дней
			latestDay = todayDay 			// обновление известного боту дня
			if diff > 0 && diff < 30 { 		// если месяц не сменялся
				posix += 86400 * diff 		// то умножаем кол-во секнд в дне на разницу дней и прибавляем к POSIX-времени
			}else{ 							// ИНАЧЕ ПИЗДЕЦ ПАНИКА, и нужен рестарт бота; ЭТО ГОЛИМЫЙ КОСТЫЛЬ И ТРЕБУЕТ ИСПРАВЛЕНИЯ
				fmt.Printf("РАЗОШЕЛСЯ DIFF")
				posix += 86400 				// На этой строке, надеемся, что просто 30-31 число сменилось на 1-ое.
			}
		}

		if obj.Message.Text == "расписос" || obj.Message.Text == "Расписос" {
			// Собираем сообщение-ответ
			b := params.NewMessagesSendBuilder()
			b.Message(getSchedule(posix))
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal("error on message send", err)
			}

		}
		if obj.Message.Text == "расписос на завтра" || obj.Message.Text == "Расписос на завтра" {
			// Собираем сообщение-ответ
			b := params.NewMessagesSendBuilder()
			b.Message(getSchedule(posix + 86400))
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal("error on message send", err)
			}
		}
	})

	err := lp.Run()
	if err != nil {
		log.Fatal(err)
	}else{
		log.Fatal("Long Poll Started")
	}
}



package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func getSchedule(unixtime int) string {

	// Объявлем переменную ответа, переменную с сылкой на расписос и переменную с "пустым" ответом
	var schedule = ""
	var scheduleUrl string = fmt.Sprintf("https://dnevnik.ru/api/userfeed/persons/1000014823656/schools/1000008291793/groups/1847068957572046690/schedule?date=%d&takeDays=1", unixtime)
	var unwanted = "{\"days\":[],\"chatStub\":{\"jid\":\"chat_students_1561239987829826940@muclight.xmpp.dnevnik.ru\"}}" // "Пустой" JSON от Дневник.Ру

	// Создаем банку с cookies и http-клиента с этими кукисами.
	options := cookiejar.Options{ // Тут открывается банка с cookies
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&options) // А тут она закрывается
	client := http.Client{Jar: jar}   // Банка теперь у клиента

	// Создаем сессию Dnevnik.ru с помощью POST-запроса
	resp, _ := client.PostForm(D_URL, url.Values{
		"login":    {A_login},
		"password": {A_password},
	})

	// Оформляем GET-запрос на расписос и записываем его в переменную data
	resp, _ = client.Get(scheduleUrl)
	data, _ := ioutil.ReadAll(resp.Body)

	// Проверяем, получили ли мы желаемый ответ от Дневник.ру
	if string(data) == unwanted {
		schedule = "Расписание не готово, или сегодня нет пар"
	} else {
		// Вбиваем JSON в структуру
		json.Unmarshal(data, &ag)
		i := 1
		for _, t := range ag.Days {
			for _, v := range t.Lessons {
				schedule += fmt.Sprint(i, "-ая пара: ", v.Subject.Name, "\nАудитория: ", v.Place, "\nВремя: ", v.Hours.StartHour, ":", v.Hours.StartMinute+"-"+v.Hours.EndHour+":"+v.Hours.EndMinute+"\n\n")
				i++
			}
		}
	}

	return schedule
}

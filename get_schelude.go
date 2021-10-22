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

func getSchedule(unixtime int) string{
	var schedule = ""
	var schedule_url string = fmt.Sprintf("https://dnevnik.ru/api/userfeed/persons/1000014823656/schools/1000008291793/groups/1847068957572046690/schedule?date=%d&takeDays=1", unixtime)
	var unwanted = "{\"days\":[],\"chatStub\":{\"jid\":\"chat_students_1561239987829826940@muclight.xmpp.dnevnik.ru\"}}" // "Пустой" JSON от Дневник.Ру

	options := cookiejar.Options{						// Тут открывается банка с cookies
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&options)					// А тут она закрывается
	client := http.Client{Jar: jar}						// Банка теперь у клиента

	resp, _ := client.PostForm(D_URL, url.Values{		// POST-запрос к странице авторизации Дневник.Ру
		"login": {A_login},								// Почему "A" - потому что А подгруппа
		"password" : {A_password},
	})

	resp, _ = client.Get(schedule_url)					// GET-запрос на страницу с JSON-расписосом
	data, _ := ioutil.ReadAll(resp.Body)				// Считывание всего содержимого GET-ответа
	//resp.Body.Close()									// Закрываем подключение
	//log.Println(string(data))   						// Лог в консоль для удобства

	if string(data) == unwanted {  								// Проверка на заполненность JSON от Дневник.ру
		schedule = "Расписание не готово, или завтра нет пар"
	}else {
		json.Unmarshal(data, &ag)						// Вбиваем JSON в структуру
		for _, t := range ag.Days {
			for _, v := range t.Lessons {
				schedule += fmt.Sprint(v.Subject.Name, ", ", v.Place, " - с ", v.Hours.StartHour, ":", v.Hours.StartMinute+" до "+v.Hours.EndHour+":"+v.Hours.EndMinute+"\n")
			}
		}
	}
	return schedule
}
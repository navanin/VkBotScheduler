package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func getSchedule(unixtime int, vkID int, peerID int) string {

	var groupID string
	var userID string
	var login string
	var password string
	var motd string

	if vkID == vk892[0] || vkID == vk892[1] {
		userID = user892
		groupID = class892
		login = login892
		password = password892
		motd = "РАСПИСАНИЕ ГРУППЫ 892"
	} else if vkID == vk782 || peerID == peer782 {
		userID = user782
		groupID = class782
		login = login782
		password = password782
		motd = "РАСПИСАНИЕ ГРУППЫ 782"
	} else if vkID == vk992 {
		userID = user992
		groupID = class992
		login = login992
		password = password992
		motd = "РАСПИСАНИЕ ГРУППЫ 992"
	}

	// Объявлем переменную ответа, переменную с сылкой на расписос и переменную с "пустым" ответом
	var schedule = ""
	var scheduleUrl string = fmt.Sprint("https://dnevnik.ru/api/userfeed/persons/", userID, "/schools/", schoolID, "/groups/", groupID, "/schedule?date=", unixtime, "&takeDays=1")
	log.Println(scheduleUrl)
	var unwanted = "\"days\":[]" // "Пустой" JSON от Дневник.Ру

	// Создаем банку с cookies и http-клиента с этими кукисами.
	options := cookiejar.Options{ // Тут открывается банка с cookies
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&options) // А тут она закрывается
	client := http.Client{Jar: jar}   // Банка теперь у клиента

	// Создаем сессию Dnevnik.ru с помощью POST-запроса
	resp, _ := client.PostForm(D_URL, url.Values{
		"login":    {login},
		"password": {password},
	})

	// Оформляем GET-запрос на расписос и записываем его в переменную data
	resp, _ = client.Get(scheduleUrl)
	data, _ := ioutil.ReadAll(resp.Body)

	// Проверяем, получили ли мы желаемый ответ от Дневник.ру
	if strings.Contains(string(data), unwanted) {
		schedule = "Расписание не готово, или сегодня нет пар"
	} else {
		schedule = fmt.Sprint(motd, "\n\n")
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

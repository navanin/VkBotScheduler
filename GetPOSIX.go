package main

import (
	"fmt"
	"time"
)

func getPOSIX() int {
	// Получаем день, месяц и год
	y, m, d := time.Now().Date()

	// Собираем из полученного дату в формате dd/mm/yyyy
	var date = fmt.Sprint(d, "/", int(m), "/", y)
	var layout = "02/01/2006"

	// Полученное время переводим в объект time.Time и возвращаем этот же объект под методом .Unix()
	tTime, _ := time.Parse(layout, date)
	return int(tTime.Unix())
}

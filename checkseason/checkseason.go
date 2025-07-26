package main

import (
	"fmt"
)

func main(){
	seasons := map[string][3]int{
		"Зима" : {1,2,12},
		"Весна" : {3,4,5},
		"Лето" : {6,7,8},
		"Осень": {9,10,11},
	}

	var month int
	fmt.Println("Введите номер месяца: ")
	fmt.Scan(&month)
	
	if (month < 1 || month > 12) {
		fmt.Println("Не корректный номер")
		return
	}

	for name, months := range seasons {
		for _,number := range months {
			if (number == month){
				fmt.Println("Сезон - ", name)
				return
			}
		}
	}
	
	fmt.Println("Сезон не найден - ")
}
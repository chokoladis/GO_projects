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
	fmt.Println("Input num month: ")
	fmt.Scan(&month)
	
	if (month < 1 || month > 12) {
		fmt.Println("Inncorrect number")
		return
	}

	for name, months := range seasons {
		for _,number := range months {
			if (number == month){
				fmt.Println("Season - ", name)
				return
			}
		}
	}
	
	fmt.Println("Season not foun - ")
}
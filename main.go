package main

import "fmt"

// pertama  = nama module

func main() {

	students := []map[string]string{
		{"name": " ruby", "score": "C"},
		{"name": "node", "score": "B"},
	}
	temp := map[string]string{
		"name":  "php",
		"score": "B",
	}
	// students.make(students[2], temp)
	// students[2] = temp
	fmt.Println(students, temp)
}

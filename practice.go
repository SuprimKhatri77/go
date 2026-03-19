package main

import "fmt"

type User struct {
	Name string
	age  int
}

func (u *User) Birthday() {
	u.age++
}

func Rename(u *User, newName string) string {
	u.Name = newName
	return u.Name
}

type Todo struct {
	Title string
}

func updateTodo(todo *Todo, newTitle *string) {
	if todo == nil {
		fmt.Println("CRITICAL: No struct provided!")
		return
	}

	if newTitle != nil {
		todo.Title = *newTitle
		if *newTitle == "" {
			fmt.Println("Note: You just set the title to be blank.")
		}
	}
}

func Practice() {

	user := &User{
		"Alice",
		25,
	}

	fmt.Println("old user: ", user)

	user.Birthday()

	fmt.Println(user.age)
	fmt.Println("New name is: ", Rename(user, "Supreme"))
	fmt.Println("new user: ", user)

	todo := &Todo{}
	fmt.Printf("todo: %v\n", todo)

	newTitle := ""

	updateTodo(todo, &newTitle)

	fmt.Println(todo.Title)
}

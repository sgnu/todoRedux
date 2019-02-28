package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/buger/goterm"
	"github.com/manifoldco/promptui"
)

//Task represents a task in the todo list
type Task struct {
	Due       Date   `json:"dueDate"`
	Category  string `json:"category"`
	Title     string `json:"title"`
	Important bool   `json:"important"`
}

//Date represents the date with month and day
type Date struct {
	Month int `json:"month"`
	Day   int `json:"day"`
}

//Returns true if d1 <= d2
func compareDates(d1, d2 Date) bool {
	if d1.Month < d2.Month {
		return true
	} else if d1.Month > d2.Month {
		return false
	} else {
		if d1.Day <= d2.Day {
			return true
		}
		return false
	}
}

func numberVal(input string) error {
	_, err := strconv.Atoi(input)
	return err
}

func main() {
	clearScreen()
	mainMenu()
}

//Acts as the main menu for todo
func mainMenu() {
	tasks := readFromFile()
	sortTasks(tasks)
	printAllTasks(tasks)
	prompt := promptui.Select{
		Label: "What would you like to do?",
		Items: []string{"New task", "Edit a task", "Complete a task", "Exit"},
	}

	_, selection, err := prompt.Run()
	check(err)

	switch selection {
	case "New task":
		clearScreen()
		addTask(tasks)
		clearScreen()
		mainMenu()
	case "Edit a task":
		clearScreen()
		editTask(tasks)
		clearScreen()
		mainMenu()
	case "Complete a task":
		clearScreen()
		completeTask(tasks)
		clearScreen()
		mainMenu()
	case "Exit":
		return
	}
}

//Asks the user for prompts to add a new task
func addTask(tasks []Task) {
	category := getUserPrompt("category")
	date := getDate()
	important := getBool("Is this important")
	title := getUserPrompt("title")

	task := Task{Due: date, Category: category, Title: title, Important: important}
	tasks = append(tasks, task)
	writeToFile(tasks)
}

//Asks the user which task and property to edit
func editTask(tasks []Task) {
	index := getTask("edit", tasks)
	if index >= len(tasks) {
		return
	}

	clearScreen()
	printTask(tasks[index])
	prompt := promptui.Select{
		Label: "Select a property to change",
		Items: []string{"Important", "Category", "Due Date", "Title"},
	}
	_, selection, _ := prompt.Run()
	switch selection {
	case "Important":
		tasks[index].Important = getBool("Is this important")
	case "Category":
		tasks[index].Category = getUserPrompt("category")
	case "Due Date":
		tasks[index].Due = getDate()
	case "Title":
		tasks[index].Title = getUserPrompt("title")
	}
	writeToFile(tasks)
}

//Asks the user which task to mark as complete and remove from the list
func completeTask(tasks []Task) {
	userInput := getTask("complete", tasks)
	if userInput >= len(tasks) {
		return
	}

	tasks = append(tasks[:userInput], tasks[userInput+1:]...)
	writeToFile(tasks)
}

//Returns a string with user input
func getUserPrompt(label string) string {
	prompt := promptui.Prompt{
		Label: "Enter the " + label,
	}

	userInput, _ := prompt.Run()
	return userInput
}

//Returns a bool with user input
func getBool(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	userInput, _ := prompt.Run()
	if userInput == "y" {
		return true
	}

	return false
}

//Returns a Date object with user input
func getDate() Date {
	monthPrompt := promptui.Select{
		Label: "Select the month",
		Items: []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
		Size:  12,
	}

	dayPrompt := promptui.Prompt{
		Label:    "Enter the day",
		Validate: numberVal,
	}

	month, _, _ := monthPrompt.Run()
	month++
	dayString, _ := dayPrompt.Run()
	day, _ := strconv.Atoi(dayString)

	return Date{Month: month, Day: day}
}

//Returns the index of a task from the list or a "Cancel" task
func getTask(label string, tasks []Task) int {
	exit := Task{Title: "Cancel"}
	templist := append(tasks, exit)

	template := &promptui.SelectTemplates{
		Active:   `{{  .Category  |  green  }}|{{  .Due.Month  |  green  }}/{{  .Due.Day  |  green  }} - {{  .Title  |  green  }}`,
		Inactive: `{{  .Category  |  blue  }}|{{  .Due.Month  |  blue  }}/{{  .Due.Day  |  blue  }} - {{  .Title  }}`,
	}

	prompt := promptui.Select{
		Label:     "Choose a task to " + label,
		Items:     templist,
		Templates: template,
	}

	index, _, _ := prompt.Run()

	return index
}

//Prints out the task list
func printAllTasks(tasks []Task) {
	for i := 0; i < len(tasks); i++ {
		important := ""
		if tasks[i].Important {
			important = goterm.Color("!*!", goterm.RED)
		}

		fmt.Printf("%3s[%10s|%02d/%02d]%40s\n", important, tasks[i].Category, tasks[i].Due.Month, tasks[i].Due.Day, tasks[i].Title)
	}
}

//Prints out a singular task
func printTask(task Task) {
	important := ""
	if task.Important {
		important = goterm.Color("!*!", goterm.RED)
	}

	fmt.Printf("%3s[%10s|%02d/%02d]%40s\n", important, task.Category, task.Due.Month, task.Due.Day, task.Title)
}

//Sorts a tasks list using insertion sort
func sortTasks(tasks []Task) {
	for i := 0; i < len(tasks); i++ {
		for j := i; j > 0 && compareDates(tasks[j].Due, tasks[j-1].Due); {
			temp := tasks[j-1]
			tasks[j-1] = tasks[j]
			tasks[j] = temp
			j--
		}
	}
}

//Checks if there is an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Gets the tasks list from the tasks file
func readFromFile() []Task {
	f, err := os.Open(os.Getenv("GOPATH") + "/tasks")
	check(err)

	defer f.Close()

	scanner := bufio.NewScanner(f)
	var arr []Task

	_ = scanner.Scan()
	jsontext := scanner.Bytes()

	_ = json.Unmarshal(jsontext, &arr)
	return arr
}

//Saves the tasks list to the tasks file
func writeToFile(tasks []Task) {
	f, err := os.Create(os.Getenv("GOPATH") + "/tasks")
	check(err)

	defer f.Close()

	data, _ := json.Marshal(tasks)

	writer := bufio.NewWriter(f)
	writer.Write(data)
	writer.Flush()
}

//Clears the screen using goterm
func clearScreen() {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()
}

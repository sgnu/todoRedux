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
	clearScreen()
	tasks := readFromFile()
	sortTasks(tasks)
	printTasks(tasks)
	prompt := promptui.Select{
		Label: "What would you like to do?",
		Items: []string{"New task", "Complete a task", "Exit"},
	}

	_, selection, err := prompt.Run()
	check(err)

	switch selection {
	case "New task":
		clearScreen()
		addTask(tasks)
		mainMenu()
	case "Complete a task":
		clearScreen()
		completeTask(tasks)
		mainMenu()
	case "Exit":
		return
	}
}

//Asks the user for prompts to add a new task
func addTask(tasks []Task) {
	var importantString string
	var important bool

	date := getDate()

	importantPrompt := promptui.Prompt{
		Label:     "Is this an important task",
		IsConfirm: true,
	}

	title := getUserPrompt("title")

	importantString, _ = importantPrompt.Run()
	if importantString == "y" {
		important = true
	} else {
		important = false
	}

	task := Task{Due: date, Title: title, Important: important}
	tasks = append(tasks, task)
	writeToFile(tasks)
}

//Asks the user which task to mark as complete and remove from the list
func completeTask(tasks []Task) {
	template := &promptui.SelectTemplates{
		Active:   `{{  .Due.Month  |  green  }}/{{  .Due.Day  |  green  }} - {{  .Title  |  green  }}`,
		Inactive: `{{  .Due.Month  |  blue  }}/{{  .Due.Day  |  blue  }} - {{  .Title  }}`,
	}

	completionPrompt := promptui.Select{
		Label:     "Which task are you completing",
		Items:     tasks,
		Templates: template,
	}

	index, _, _ := completionPrompt.Run()
	tasks = append(tasks[:index], tasks[index+1:]...)
	writeToFile(tasks)
}

func getUserPrompt(label string) string {
	prompt := promptui.Prompt{
		Label: "Enter the " + label,
	}

	userInput, _ := prompt.Run()
	return userInput
}

//Creates a prompt asking for a month and day
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

//Prints out the task list
func printTasks(tasks []Task) {
	for i := 0; i < len(tasks); i++ {
		important := ""
		if tasks[i].Important {
			important = goterm.Color("!*!", goterm.RED)
		}

		fmt.Printf("%3s[%10s|%02d/%02d]%40s\n", important, tasks[i].Category, tasks[i].Due.Month, tasks[i].Due.Day, tasks[i].Title)
	}
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

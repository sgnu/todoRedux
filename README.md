# todoRedux

A CLI todo list app with somewhat of an emphasis on *a e s t h e t i c s*. A successor to my [todo app](https://github.com/sgnu/todo).
![screenshot](screenshot.png)

## Installation

todoRedux can be installed using the command:

```sh
go get github.com/sgnu/todoRedux
```

You will need to create a `tasks` file (with no extension) in your `$GOPATH` directory. It can be empty, but the file must exist.

## Usage

Using this app is pretty self explanatory:

- Add a task
  - Adds a new task to the list
- Edit a task
  - Edits a task that is on the list
- Complete a task
  - Completes and removes a task that is on the list
- Exit
  - Closes the app

## How it works

Each task is represented by a Task struct which has the fields:

```go
Due       Date    // When the task is due
Category  string  // The category of the task
Title     string  // The name of the task
Important bool    // If the task is important
```

And the Date struct has the fields:

```go
Month int // Month (1 to 12)
Day   int // Day (1 to 31)
```

The majority of the UI is handled by [promptui](https://github.com/manifoldco/promptui) and [goterm](https://github.com/buger/goterm) for clearing the screen.

A slice containing task structs, `tasks[]`, is used throughout. It is written to and read from the `tasks` file by converting to and from JSON.

Method names describe themselves and also include a short description commented directly before them.
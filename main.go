package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type TodoItem struct {
	text     string
	isDone   bool
}

var todoList []TodoItem
var listMutex = &sync.Mutex{}

func main() {
	ln, _ := net.Listen("tcp", ":8888")

	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		listMutex.Lock()
		handleMessage(scanner.Text(), conn)
		listMutex.Unlock()
	}
}

func handleMessage(msg string, conn net.Conn) {
	parts := strings.SplitN(msg, " ", 2)
	command := strings.ToUpper(parts[0])

	switch command {
	case "SHOW":
		for i, item := range todoList {
			doneStr := " "
			if item.isDone {
				doneStr = "x"
			}
			fmt.Fprintf(conn, "%d [%s] %s\n", i, doneStr, item.text)
		}
	case "UPDATE":
		args := strings.SplitN(parts[1], " ", 2)
		index, _ := strconv.Atoi(args[0])
		todoList[index].text = args[1]
	case "ADD":
		todoList = append(todoList, TodoItem{parts[1], false})
	case "DELETE":
		index, _ := strconv.Atoi(parts[1])
		todoList = append(todoList[:index], todoList[index+1:]...)
	case "CHECK":
		index, _ := strconv.Atoi(parts[1])
		todoList[index].isDone = true
	case "UNCHECK":
		index, _ := strconv.Atoi(parts[1])
		todoList[index].isDone = false
	case "HELP":
		fmt.Fprintf(conn, "Hi! Available commands are:\nSHOW - display the list\nUPDATE <number> <text> - Update list item at <number> with <text>\nADD <text> - Add an item to the list\nDELETE <number> - Delete the item at <number>\nCHECK <number> - Mark item at <number> as complete\nUNCHECK <number> - Mark item at <number> as incomplete\nHELP - Show this help message\n")
	default:
		fmt.Fprintf(conn, "Unknown command: %s\n", command)
	}
}
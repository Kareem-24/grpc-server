package main

import (
	"ToDoList/pb"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := "localhost:8081"
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewTodoServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Get a todo")
		fmt.Println("2. Get all todos")
		fmt.Println("3. Create new todo")
		fmt.Println("4. Update a todo")
		fmt.Println("5. Delete a todo")
		fmt.Println("6. Exit")

		choice, err := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		if err != nil {
			fmt.Println("Error in:", err)
		}
		switch choice {
		case "1": // get todo by id
			fmt.Print("enter todo id :")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			req := &pb.TodoId{Id: int32(id)}

			res, err := client.GetTodo(context.Background(), req)

			if err != nil {
				fmt.Println("Error getting todo:", err)
			} else {
				fmt.Println("Todo:", res.GetTodo())
			}
		case "2": // get all todos
			req := &pb.GetTodos{}
			res, err := client.GetAllTodo(context.Background(), req)

			if err != nil {
				fmt.Println("Error getting todo:", err)
			} else {
				fmt.Println("Todo:", res.GetTodo())
			}

		case "3": // create new todo
			fmt.Print("enter todo title :")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("enter todo description :")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)

			updateReq := &pb.Todo{
				Title:       title,
				Description: desc,
			}

			res, err := client.CreateTodo(context.Background(), updateReq)
			if err != nil {
				fmt.Println("Error creating todo:", err)
			} else {
				fmt.Println("Todo:", res.GetTodo())
			}

		case "4": // update a todo
			fmt.Print("enter todo id :")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			req := &pb.TodoId{Id: int32(id)}

			res, err := client.GetTodo(context.Background(), req)

			if err != nil {
				fmt.Println("Error getting todo:", err)
				return
			} else {
				fmt.Println("Todo:", res.GetTodo())
			}

			fmt.Print("enter todo title :")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("enter todo description :")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)

			updateReq := &pb.Todo{
				Id:          int32(id),
				Title:       title,
				Description: desc,
			}

			res, err = client.UpdateTodo(context.Background(), updateReq)
			if err != nil {
				fmt.Println("Error updating todo:", err)
			} else {
				fmt.Println("Todo:", res.GetTodo())
			}

		case "5":
			fmt.Print("enter todo id :")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			req := &pb.TodoId{Id: int32(id)}

			res, err := client.DeleteTodo(context.Background(), req)
			if err != nil {
				fmt.Println("Error getting todo:", err)
			} else {
				fmt.Println("Todo:", res.GetTodo())
			}
		case "6":
			fmt.Println("Exiting...")
			return
		default:
			data, _ := reader.ReadString('\n')

			// fmt.Println("Invalid choice")
			fmt.Printf("Invalid choice %v", data)
		}

	}

}

package main

import (
	"ToDoList/entities"
	"ToDoList/pb"
	service "ToDoList/server/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	db := entities.InitDB()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to connect : %v", err)
	}

	serv := grpc.NewServer()
	todoServ := &service.TodoServiceServer{Db: db}

	pb.RegisterTodoServiceServer(serv, todoServ)
	log.Printf("server listining on %v", lis.Addr())
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("failed to connect : %v", err)
	}
}

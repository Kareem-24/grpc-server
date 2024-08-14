package main

import (
	"ToDoList/entities"
	"ToDoList/pb"
	ch "ToDoList/serverWithChannles/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	db := entities.InitDB()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to connect : %v", err)
	}

	serv := grpc.NewServer()
	todoServ := ch.NewTodoServiceServer(db)

	pb.RegisterTodoServiceServer(serv, todoServ)
	log.Printf("server listining on %v", lis.Addr())
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("failed to connect : %v", err)
	}
}

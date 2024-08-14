package services

import (
	"ToDoList/entities"
	"ToDoList/pb"
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type TodoServiceServer struct {
	Db *gorm.DB
	pb.UnimplementedTodoServiceServer
	mu sync.Mutex
}

func (s *TodoServiceServer) GetTodo(ctx context.Context, req *pb.TodoId) (*pb.TodoResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid todo id")
	}

	var todo entities.Todo
	if err := s.Db.First(&todo, req.GetId()).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "todo not found: %v", err)
	}

	response := &pb.TodoResponse{
		Todo: &pb.Todo{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
			IsCompleted: todo.IsCompleted,
		},
	}
	return response, nil
}

func (s *TodoServiceServer) GetAllTodo(ctx context.Context, req *pb.GetTodos) (*pb.TodoList, error) {
	var todolist []entities.Todo
	res := s.Db.Find(&todolist)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todoList: %v ", res.Error)

	}

	pbTodos := make([]*pb.Todo, len(todolist))
	for i, todo := range todolist {
		pbTodos[i] = &pb.Todo{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
		}
	}

	return &pb.TodoList{Todo: pbTodos}, nil
}

func (s *TodoServiceServer) CreateTodo(ctx context.Context, req *pb.Todo) (*pb.TodoList, error) {
	if req.GetDescription() == "" && req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title or description is empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	res := s.Db.Create(&req)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todoList: %v ", res.Error)
	}
	return s.GetAllTodo(ctx, &pb.GetTodos{})
}

func (s *TodoServiceServer) UpdateTodo(ctx context.Context, req *pb.Todo) (*pb.TodoResponse, error) {

	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid todo id")
	}

	if req.GetDescription() == "" && req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title or description is empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	var todo entities.Todo
	if err := s.Db.First(&todo, req.GetId()).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "todo not found: %v", err)
	}

	todo.Title = req.GetTitle()
	todo.Description = req.GetDescription()
	todo.IsCompleted = req.GetIsCompleted()

	if err := s.Db.Save(&todo).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update todo: %v", err)
	}

	response := &pb.TodoResponse{
		Todo: &pb.Todo{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
			IsCompleted: todo.IsCompleted,
		},
	}

	return response, nil
}

func (s *TodoServiceServer) DeleteTodo(ctx context.Context, req *pb.TodoId) (*pb.TodoList, error) {

	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid todo id")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	var todo entities.Todo

	if err := s.Db.Delete(&todo, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete todo: %v", err)
	}
	return s.GetAllTodo(ctx, &pb.GetTodos{})

}

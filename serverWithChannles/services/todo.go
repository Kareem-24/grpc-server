package serverWithChannles

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
	mu             sync.Mutex
	getTodoChan    chan GetTodoRequest
	getAllTodoChan chan GetAllTodoRequest
	createTodoChan chan CreateTodoRequest
	updateTodoChan chan UpdateTodoRequest
	deleteTodoChan chan DeleteTodoRequest
}

type GetTodoRequest struct {
	ctx      context.Context
	req      *pb.TodoId
	response chan<- *pb.TodoResponse
	err      chan<- error
}

type GetAllTodoRequest struct {
	ctx      context.Context
	req      *pb.GetTodos
	response chan<- *pb.TodoList
	err      chan<- error
}

type CreateTodoRequest struct {
	ctx      context.Context
	req      *pb.Todo
	response chan<- *pb.TodoList
	err      chan<- error
}

type UpdateTodoRequest struct {
	ctx      context.Context
	req      *pb.Todo
	response chan<- *pb.TodoResponse
	err      chan<- error
}

type DeleteTodoRequest struct {
	ctx      context.Context
	req      *pb.TodoId
	response chan<- *pb.TodoList
	err      chan<- error
}

func NewTodoServiceServer(db *gorm.DB) *TodoServiceServer {
	s := &TodoServiceServer{
		Db:             db,
		getTodoChan:    make(chan GetTodoRequest),
		getAllTodoChan: make(chan GetAllTodoRequest),
		createTodoChan: make(chan CreateTodoRequest),
		updateTodoChan: make(chan UpdateTodoRequest),
		deleteTodoChan: make(chan DeleteTodoRequest),
	}

	go s.handleGetTodoRequests()
	go s.handleGetAllTodoRequests()
	go s.handleCreateTodoRequests()
	go s.handleUpdateTodoRequests()
	go s.handleDeleteTodoRequests()

	return s
}

func (s *TodoServiceServer) GetTodo(ctx context.Context, req *pb.TodoId) (*pb.TodoResponse, error) {
	response := make(chan *pb.TodoResponse)
	err := make(chan error)
	s.getTodoChan <- GetTodoRequest{ctx, req, response, err}
	return <-response, <-err
}

func (s *TodoServiceServer) GetAllTodo(ctx context.Context, req *pb.GetTodos) (*pb.TodoList, error) {
	response := make(chan *pb.TodoList)
	err := make(chan error)
	s.getAllTodoChan <- GetAllTodoRequest{ctx, req, response, err}
	return <-response, <-err
}

func (s *TodoServiceServer) CreateTodo(ctx context.Context, req *pb.Todo) (*pb.TodoList, error) {
	response := make(chan *pb.TodoList)
	err := make(chan error)
	s.createTodoChan <- CreateTodoRequest{ctx, req, response, err}
	return <-response, <-err
}

func (s *TodoServiceServer) UpdateTodo(ctx context.Context, req *pb.Todo) (*pb.TodoResponse, error) {
	response := make(chan *pb.TodoResponse)
	err := make(chan error)
	s.updateTodoChan <- UpdateTodoRequest{ctx, req, response, err}
	return <-response, <-err
}

func (s *TodoServiceServer) DeleteTodo(ctx context.Context, req *pb.TodoId) (*pb.TodoList, error) {
	response := make(chan *pb.TodoList)
	err := make(chan error)
	s.deleteTodoChan <- DeleteTodoRequest{ctx, req, response, err}
	return <-response, <-err
}

func (s *TodoServiceServer) handleGetTodoRequests() {
	for req := range s.getTodoChan {
		var todo entities.Todo
		if err := s.Db.First(&todo, req.req.GetId()).Error; err != nil {
			req.response <- nil
			req.err <- status.Errorf(codes.NotFound, "todo not found: %v", err)
		} else {
			req.response <- &pb.TodoResponse{
				Todo: &pb.Todo{
					Id:          todo.Id,
					Title:       todo.Title,
					Description: todo.Description,
					IsCompleted: todo.IsCompleted,
				},
			}
			req.err <- nil
		}
	}
}
func (s *TodoServiceServer) handleGetAllTodoRequests() {
	for req := range s.getAllTodoChan {
		var todolist []entities.Todo
		if res := s.Db.Find(&todolist); res.Error != nil {
			req.response <- nil
			req.err <- status.Errorf(codes.Internal, "failed to get todoList: %v ", res.Error)
		} else {
			pbTodos := make([]*pb.Todo, len(todolist))
			for i, todo := range todolist {
				pbTodos[i] = &pb.Todo{
					Id:          todo.Id,
					Title:       todo.Title,
					Description: todo.Description,
					IsCompleted: todo.IsCompleted,
				}
			}
			req.response <- &pb.TodoList{Todo: pbTodos}
			req.err <- nil
		}
	}
}

func (s *TodoServiceServer) handleCreateTodoRequests() {
	for req := range s.createTodoChan {
		s.mu.Lock()
		if res := s.Db.Create(&req.req); res.Error != nil {
			req.response <- nil
			req.err <- status.Errorf(codes.Internal, "failed to create todo: %v", res.Error)
		} else {
			todos, _ := s.GetAllTodo(req.ctx, &pb.GetTodos{})
			req.response <- todos
			req.err <- nil
		}
		s.mu.Unlock()
	}
}

func (s *TodoServiceServer) handleUpdateTodoRequests() {
	for req := range s.updateTodoChan {
		s.mu.Lock()
		var todo entities.Todo
		if err := s.Db.First(&todo, req.req.GetId()).Error; err != nil {
			req.response <- nil
			req.err <- status.Errorf(codes.NotFound, "todo not found: %v", err)
		} else {
			todo.Title = req.req.GetTitle()
			todo.Description = req.req.GetDescription()
			todo.IsCompleted = req.req.GetIsCompleted()

			if err := s.Db.Save(&todo).Error; err != nil {
				req.response <- nil
				req.err <- status.Errorf(codes.Internal, "failed to update todo: %v", err)
			} else {
				req.response <- &pb.TodoResponse{
					Todo: &pb.Todo{
						Id:          todo.Id,
						Title:       todo.Title,
						Description: todo.Description,
						IsCompleted: todo.IsCompleted,
					},
				}
				req.err <- nil
			}
		}
		s.mu.Unlock()
	}
}

func (s *TodoServiceServer) handleDeleteTodoRequests() {
	for req := range s.deleteTodoChan {
		s.mu.Lock()
		if err := s.Db.Delete(&entities.Todo{}, req.req.GetId()).Error; err != nil {
			req.response <- nil
			req.err <- status.Errorf(codes.Internal, "failed to delete todo: %v", err)
		} else {
			todos, _ := s.GetAllTodo(req.ctx, &pb.GetTodos{})
			req.response <- todos
			req.err <- nil
		}
		s.mu.Unlock()
	}
}

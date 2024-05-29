package main

import (
	"context"
	"fmt"
	"grpc-redis-server/protos/todo/protos/todo"
	"math/rand"

	"github.com/go-redis/redis"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/proto"
)

type TodoServiceServer struct {
	redisClient *redis.Client
	log         hclog.Logger
	todo.UnimplementedTodoServiceServer
}

func NewTodoServiceServer(redisAddr string, log hclog.Logger) *TodoServiceServer {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &TodoServiceServer{
		redisClient: client,
		log:         log,
	}
}

func generateUniqueID() string {
	min := 1000
	max := 9999
	return fmt.Sprintf("%d", rand.Intn(max-min+1)+min)
}

func (s *TodoServiceServer) AddTodo(ctx context.Context, addReq *todo.AddTodoRequest) (*todo.AddTodoResponse, error) {
	id := generateUniqueID()
	todoItem := &todo.TodoItem{
		Id:          id,
		Title:       addReq.Title,
		Description: addReq.Description,
	}

	todoItemBytes, err := proto.Marshal(todoItem)
	if err != nil {
		return nil, err
	}

	err = s.redisClient.Set(id, todoItemBytes, 0).Err()
	if err != nil {
		return nil, err
	}

	s.log.Info("AddTodo", "ID", id, "Title", addReq.Title, "Description", addReq.Description)

	return &todo.AddTodoResponse{Id: id}, nil
}

func (s *TodoServiceServer) GetTodo(ctx context.Context, getReq *todo.GetTodoRequest) (*todo.GetTodoResponse, error) {
	val, err := s.redisClient.Get(getReq.Id).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("todo item not found")
	} else if err != nil {
		return nil, err
	}

	var todoItem todo.TodoItem
	err = proto.Unmarshal([]byte(val), &todoItem)
	if err != nil {
		return nil, err
	}

	s.log.Info("GetTodo", "ID", todoItem.Id)

	return &todo.GetTodoResponse{
		Id:          todoItem.Id,
		Title:       todoItem.Title,
		Description: todoItem.Description,
	}, nil
}

func (s *TodoServiceServer) UpdateTodo(ctx context.Context, updateReq *todo.UpdateTodoRequest) (*todo.UpdateTodoResponse, error) {
	val, err := s.redisClient.Get(updateReq.Id).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("todo item not found")
	} else if err != nil {
		return nil, err
	}

	var todoItem todo.TodoItem
	err = proto.Unmarshal([]byte(val), &todoItem)
	if err != nil {
		return nil, err
	}

	todoItem.Title = updateReq.Title
	todoItem.Description = updateReq.Description

	todoItemBytes, err := proto.Marshal(&todoItem)
	if err != nil {
		return nil, err
	}

	err = s.redisClient.Set(updateReq.Id, todoItemBytes, 0).Err()
	if err != nil {
		return nil, err
	}

	s.log.Info("UpdateTodo", "ID", updateReq.Id, "Title", updateReq.Title, "Description", updateReq.Description)

	return &todo.UpdateTodoResponse{Success: "true"}, nil
}

func (s *TodoServiceServer) DeleteTodo(ctx context.Context, delReq *todo.DeleteTodoRequest) (*todo.DeleteTodoResponse, error) {
	err := s.redisClient.Del(delReq.Id).Err()
	if err != nil {
		return nil, err
	}

	s.log.Info("DeleteTodo", "ID", delReq.Id)

	return &todo.DeleteTodoResponse{Success: "true"}, nil
}

func (s *TodoServiceServer) ListTodos(ctx context.Context, listReq *todo.ListTodosRequest) (*todo.ListTodosResponse, error) {
	var cursor uint64
	var todoItems []*todo.TodoItem

	for {
		keys, nextCursor, err := s.redisClient.Scan(cursor, "*", 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			val, err := s.redisClient.Get(key).Result()
			if err != nil {
				return nil, err
			}

			var todoItem todo.TodoItem
			err = proto.Unmarshal([]byte(val), &todoItem)
			if err != nil {
				return nil, err
			}

			todoItems = append(todoItems, &todoItem)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	s.log.Info("ListTodos")

	return &todo.ListTodosResponse{Todos: todoItems}, nil
}

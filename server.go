package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	api "github.com/ystkfujii/example-oapi-codegen/openapi"
)

type Server struct {
	users    map[int]User
	nextID   int
	usersMux sync.RWMutex
}

var _ api.ServerInterface = (*Server)(nil)

func NewServer() *Server {
	return &Server{
		users:  make(map[int]User),
		nextID: 1,
	}
}

func (s *Server) GetUsers(ctx echo.Context) error {
	fmt.Println("GET /users called")
	s.usersMux.RLock()
	defer s.usersMux.RUnlock()

	users := make([]api.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, userConv.ToOpenAPIUser(user))
	}

	return ctx.JSON(http.StatusOK, users)
}

func (s *Server) PostUsers(ctx echo.Context) error {
	fmt.Println("POST /users called")

	var newUser api.NewUser
	if err := ctx.Bind(&newUser); err != nil {
		return ReturnInvalidRequestBodyError(ctx, err)
	}

	user, err := userConv.ToDomainUser(api.User{
		Id:   s.nextID,
		Name: newUser.Name,
		Age:  newUser.Age,
	})
	if err != nil {
		return ReturnInvalidUserDataError(ctx, err)
	}

	s.usersMux.Lock()
	defer s.usersMux.Unlock()

	s.users[s.nextID] = user
	s.nextID++

	return ctx.JSON(http.StatusCreated, userConv.ToOpenAPIUser(user))
}

func (s *Server) GetUsersId(ctx echo.Context, id int) error {
	fmt.Printf("GET /users/%d called\n", id)
	s.usersMux.RLock()
	defer s.usersMux.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return ReturnUserNotFoundError(ctx, id)
	}

	return ctx.JSON(http.StatusOK, userConv.ToOpenAPIUser(user))
}

func (s *Server) DeleteUsersId(ctx echo.Context, id int) error {
	fmt.Printf("DELETE /users/%d called\n", id)
	s.usersMux.Lock()
	defer s.usersMux.Unlock()

	if _, exists := s.users[id]; !exists {
		return ReturnUserNotFoundError(ctx, id)
	}

	delete(s.users, id)
	return ctx.NoContent(http.StatusNoContent)
}

var userConv _userConv

type _userConv struct{}

func (c _userConv) ToOpenAPIUser(u User) api.User {
	return api.User{
		Id:   u.Id,
		Name: api.Name{First: u.Name.First, Last: u.Name.Last, Middle: u.Name.Middle},
		Age:  u.Age,
	}
}

func (c _userConv) ToDomainUser(oapiUser api.User) (User, error) {

	name, err := NewName(oapiUser.Name.First, oapiUser.Name.Last, oapiUser.Name.Middle)
	if err != nil {
		return User{}, err
	}
	user, err := NewUser(oapiUser.Age, oapiUser.Id, name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

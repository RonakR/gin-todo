package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

type Store struct {
	mu    sync.Mutex
	next  int
	todos map[int]Todo
}

func NewStore() *Store {
	return &Store{
		next:  1,
		todos: make(map[int]Todo),
	}
}

func (s *Store) List() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		out = append(out, todo)
	}
	return out
}

func (s *Store) Get(id int) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	return todo, ok
}

func (s *Store) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := Todo{
		ID:        s.next,
		Title:     title,
		Completed: false,
	}
	s.todos[todo.ID] = todo
	s.next++
	return todo
}

func (s *Store) Update(id int, req UpdateTodoRequest) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return Todo{}, false
	}
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	s.todos[id] = todo
	return todo, true
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return false
	}
	delete(s.todos, id)
	return true
}

func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

func main() {
	store := NewStore()
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, store.List())
	})

	router.POST("/todos", func(c *gin.Context) {
		var req CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
			return
		}
		c.JSON(http.StatusCreated, store.Create(req.Title))
	})

	router.GET("/todos/:id", func(c *gin.Context) {
		id, ok := parseID(c)
		if !ok {
			return
		}
		todo, found := store.Get(id)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusOK, todo)
	})

	router.PATCH("/todos/:id", func(c *gin.Context) {
		id, ok := parseID(c)
		if !ok {
			return
		}
		var req UpdateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		if req.Title == nil && req.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}
		todo, found := store.Update(id, req)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusOK, todo)
	})

	router.DELETE("/todos/:id", func(c *gin.Context) {
		id, ok := parseID(c)
		if !ok {
			return
		}
		if !store.Delete(id) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	_ = router.Run(":8080")
}

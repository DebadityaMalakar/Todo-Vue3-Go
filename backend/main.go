package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Task structure to store task data
type Task struct {
	ID          int       `json:"id"`
	Task        string    `json:"task"`
	CreatedAt   time.Time `json:"createdAt"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	Completed   bool      `json:"completed"`
}

// Global slice to store tasks
var tasks []Task

// WebSocket connections
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []Task)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for WebSocket
		return true
	},
}

// Handler to broadcast tasks to all connected clients
func broadcastTasks() {
	broadcast <- tasks
}

// Handler for the /tasks WebSocket route
func getTasksWS(c *gin.Context) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket:", err)
		return
	}
	defer func() {
		conn.Close()
		delete(clients, conn)
	}()

	// Register new client
	clients[conn] = true

	// Send initial tasks
	err = conn.WriteJSON(tasks)
	if err != nil {
		log.Println("Error sending initial tasks:", err)
		return
	}

	// Keep connection open and listen for messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// Broadcast loop
func handleBroadcast() {
	for {
		// Broadcast tasks to all connected clients
		tasksToBroadcast := <-broadcast
		for client := range clients {
			err := client.WriteJSON(tasksToBroadcast)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// Handler for the /tasks HTTP route (GET)
func getTasks(c *gin.Context) {
	// Return the tasks list (empty if no tasks)
	c.JSON(http.StatusOK, tasks)
}

// Handler for the /tasks HTTP route (POST)
func createTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simulate adding the new task
	newTask.ID = len(tasks) + 1
	newTask.CreatedAt = time.Now()
	tasks = append(tasks, newTask)

	// Broadcast updated tasks
	broadcastTasks()

	// Respond with the status
	c.JSON(http.StatusOK, gin.H{"status": "Task created successfully"})
}

// Handler to mark a task as completed
func completeTask(c *gin.Context) {
	// Get task ID from URL parameter
	idStr := c.Param("id")

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Find the task
	for i, task := range tasks {
		if task.ID == id {
			// Mark task as completed
			tasks[i].Completed = true
			tasks[i].CompletedAt = time.Now()

			// Broadcast updated tasks
			broadcastTasks()

			c.JSON(http.StatusOK, gin.H{"status": "Task marked as completed"})
			return
		}
	}

	// Task not found
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

// Handler for deleting a task
func deleteTask(c *gin.Context) {
	// Get task ID from URL parameter
	idStr := c.Param("id")

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Find and delete the task
	for i, task := range tasks {
		if task.ID == id {
			// Remove task from the slice
			tasks = append(tasks[:i], tasks[i+1:]...)

			// Broadcast updated tasks
			broadcastTasks()

			c.JSON(http.StatusOK, gin.H{"status": "Task deleted successfully"})
			return
		}
	}

	// Task not found
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func main() {
	// Set up Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Vue.js dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Start broadcast handler
	go handleBroadcast()

	// Routes for REST API
	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id/complete", completeTask)
	r.DELETE("/tasks/:id", deleteTask)

	// Route for WebSocket
	r.GET("/tasks/ws", getTasksWS)

	// Run the server
	r.Run(":8080")
}

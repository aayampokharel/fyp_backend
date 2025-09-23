package main

import (
	"context"
	"fmt"
	"net/http"
	middleware "project/internals/domain/middlewares"
	logger "project/package/utils/pkg"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

// Book represents a book in our library
type Book struct {
	ID     string `json:"id" doc:"Unique book identifier" example:"1"`
	Title  string `json:"title" doc:"Book title" maxLength:"200" example:"The Go Programming Language"`
	Author string `json:"author" doc:"Book author" maxLength:"100" example:"Alan A. A. Donovan"`
	Year   int    `json:"year" doc:"Publication year" minimum:"1000" maximum:"2024" example:"2015"`
}

// BookInput for creating/updating books
type BookInput struct {
	Body *Book `doc:"Book data"`
}

// BookIDInput for operations requiring a book ID
type BookIDInput struct {
	ID string `path:"id" doc:"Book ID" example:"1"`
}

// BooksResponse for listing books
type BooksResponse struct {
	Books []Book `json:"books" doc:"List of books"`
}

// In-memory storage
var books = map[string]Book{
	"1": {ID: "1", Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Year: 2015},
	"2": {ID: "2", Title: "Clean Code", Author: "Robert C. Martin", Year: 2008},
}

func main() {
	// Create Gin router
	router := gin.New()
	router.Use(middleware.ColorStatusLogger())

	// Create Huma config with docs enabled
	config := huma.DefaultConfig("Book Library API", "1.0.0")
	config.DocsPath = "/docs" // Enable docs endpoint
	config.OpenAPI.Info.Description = "A simple book management API demonstrating Huma framework"

	// Create Huma API
	// Create Huma API - FIXED
	config = huma.DefaultConfig("Book Library API", "1.0.0")
	config.DocsPath = "/docs"
	config.OpenAPI.Info.Description = "A simple book management API demonstrating Huma framework"

	// Use humagin.New: pass router + config
	api := humagin.New(router, config)

	// 1. Create a new book
	huma.Register(api, huma.Operation{
		OperationID: "create-book",
		Method:      http.MethodPost,
		Path:        "/books",
		Summary:     "Create a new book",
		Tags:        []string{"Books"},
	}, func(ctx context.Context, input *BookInput) (*Book, error) {
		book := input.Body
		if book.ID == "" {
			book.ID = fmt.Sprintf("%d", len(books)+1)
		}
		books[book.ID] = *book
		return book, nil
	})

	// 2. Get all books
	huma.Register(api, huma.Operation{
		OperationID: "list-books",
		Method:      http.MethodGet,
		Path:        "/books",
		Summary:     "Get all books",
		Tags:        []string{"Books"},
	}, func(ctx context.Context, input *struct{}) (*BooksResponse, error) {
		bookList := make([]Book, 0, len(books))
		for _, book := range books {
			bookList = append(bookList, book)
		}
		return &BooksResponse{Books: bookList}, nil
	})

	// 3. Get a specific book by ID
	huma.Register(api, huma.Operation{
		OperationID: "get-book",
		Method:      http.MethodGet,
		Path:        "/books/{id}",
		Summary:     "Get a book by ID",
		Tags:        []string{"Books"},
	}, func(ctx context.Context, input *BookIDInput) (*Book, error) {
		book, exists := books[input.ID]
		if !exists {
			return nil, huma.Error404NotFound("Book not found")
		}
		return &book, nil
	})

	// 4. Update a book
	huma.Register(api, huma.Operation{
		OperationID: "update-book",
		Method:      http.MethodPut,
		Path:        "/books/{id}",
		Summary:     "Update a book",
		Tags:        []string{"Books"},
	}, func(ctx context.Context, input *BookInput) (*Book, error) {
		book := input.Body
		if _, exists := books[book.ID]; !exists {
			return nil, huma.Error404NotFound("Book not found")
		}
		books[book.ID] = *book
		return book, nil
	})

	// 5. Delete a book
	huma.Register(api, huma.Operation{
		OperationID: "delete-book",
		Method:      http.MethodDelete,
		Path:        "/books/{id}",
		Summary:     "Delete a book",
		Tags:        []string{"Books"},
	}, func(ctx context.Context, input *BookIDInput) (*struct{}, error) {
		if _, exists := books[input.ID]; !exists {
			return nil, huma.Error404NotFound("Book not found")
		}
		delete(books, input.ID)
		return &struct{}{}, nil
	})

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "books_count": len(books)})
	})
	logger.InitLogger()



	logger.InitLogger()
	logger.Logger.Infoln("🌱 Starting the app...")

	fmt.Println("📚 Book Library API Server starting...")
	fmt.Println("📍 Base URL: http://localhost:8888")
	fmt.Println("📖 API Docs: http://localhost:8888/docs")
	fmt.Println("🔍 OpenAPI Spec: http://localhost:8888/openapi.json")
	fmt.Println("❤️  Health Check: http://localhost:8888/health")
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  POST   /books      - Create a new book")
	fmt.Println("  GET    /books      - List all books")
	fmt.Println("  GET    /books/{id} - Get a specific book")
	fmt.Println("  PUT    /books/{id} - Update a book")
	fmt.Println("  DELETE /books/{id} - Delete a book")

	// Start server
	if err := router.Run(":8888"); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

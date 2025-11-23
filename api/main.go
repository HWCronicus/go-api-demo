package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"

	"go-api-demo/database"
	_ "go-api-demo/docs"
)

var queries *database.Queries
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uuid.UUID `json:"id"`
		Email string    `json:"email"`
	} `json:"user"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type DeleteCommentRequest struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type DeleteCommentSuccessResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// @title Go Demo API
// @version 1.0
// @description Go Demo API with user authentication and comments

// @contact.name API Support
// @contact.email alan.e.george86@gmail.com

// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := validateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func initDB() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	queries = database.New(db)
	return nil
}

func generateToken(userID uuid.UUID, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./html")).ServeHTTP(w, r)
}

func handleResume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=Alan_George_Resume.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, "files/Alan_George_Resume.pdf")
}

// handleCreateUser creates a new user
// @Summary Create a new user
// @Description Register a new user with name and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User registration details"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user [post]
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Name and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	user, err := queries.CreateUser(ctx, database.CreateUserParams{
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		http.Error(w, "Error creating user (name may already exist)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	})
}

// handleLogin authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticate with username and password to receive a JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Name and password are required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	user, err := queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: "Bearer " + token,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleCreateComment creates a new comment (requires authentication)
// @Summary Create a comment
// @Description Post a new comment (requires JWT authentication)
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment body CreateCommentRequest true "Comment content"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comment [post]
func handleCreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	comment, err := queries.CreateComment(ctx, database.CreateCommentParams{
		UserID:  uuid.NullUUID{UUID: claims.UserID, Valid: true},
		Content: req.Content,
	})

	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommentResponse{
		ID:        comment.ID,
		Email:     comment.Email,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Time,
	})
}

// handleListComments lists all comments
// @Summary List all comments
// @Description List all comments
// @Tags comments
// @Produce json
// @Success 200 {object} []CommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments [get]
func handleListComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := context.Background()
	comments, err := queries.ListComments(ctx)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}
	var response []CommentResponse
	for _, comment := range comments {
		response = append(response, CommentResponse{
			ID:        comment.ID,
			Email:     comment.Email,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.Time,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleDeleteComment deletes a comment (requires authentication)
// @Summary Delete a comment
// @Description Delete a comment (requires JWT authentication)
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment body DeleteCommentRequest true "Comment deletion details"
// @Success 200 {object} DeleteCommentSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comment [delete]
func handleDeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()

	var req DeleteCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == uuid.Nil {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	user, getUserError := queries.GetUserByEmail(ctx, req.Email)

	if getUserError != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.ID != claims.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	deleteError := queries.DeleteComment(ctx, database.DeleteCommentParams{
		ID:    req.ID,
		Email: req.Email,
	})

	if deleteError != nil {
		http.Error(w, "Error deleting comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeleteCommentSuccessResponse{
		Message: "Comment deleted successfully",
	})
}

func main() {
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET environment variable must be set")
	}

	time.Sleep(1 * time.Second)

	if err := initDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/resume", handleResume)
	mux.HandleFunc("/user", handleCreateUser)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("POST /comment", authMiddleware(handleCreateComment))
	mux.HandleFunc("DELETE /comment", authMiddleware(handleDeleteComment))
	mux.HandleFunc("/comments", handleListComments)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server starting...")

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMiddleware(mux),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}

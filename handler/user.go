package handler

import (
	"be-education/dto"
	"be-education/models"
	"be-education/service"
	"be-education/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandlerImpl struct {
	userService service.UserService
	baseURL     string
}

func NewUserHandler(userService service.UserService, baseURL string) userHandlerImpl {
	return userHandlerImpl{userService: userService, baseURL: baseURL}
}

func (h *userHandlerImpl) CreateMahasiswa(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or validation failed", "details": err.Error()})
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Class:    req.Class,
		Birthday: req.Birthday,
		Role:     "mahasiswa",
	}

	err := h.userService.CreateUser(c.Request.Context(), user)
	if err != nil {
		switch err.Error() {
		case fmt.Sprintf("user with email %s already exists", user.Email):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case "email cannot be empty":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "password cannot be empty for hashing":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *userHandlerImpl) Login(c *gin.Context) {
	var req dto.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or validation failed", "details": err.Error()})
		return
	}

	token, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func (h *userHandlerImpl) GetProfile(c *gin.Context) {
	claims, exists := utils.GetCurrentUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User claims not found in context. Authentication required."})
		return
	}

	userID := claims.UserID

	ctx := c.Request.Context()

	userDTO, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		fmt.Printf("Error getting user from service for ID %d: %v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile due to internal error"})
		return
	}

	if userDTO == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User profile with ID %d not found", userID)})
		return
	}

	c.JSON(http.StatusOK, userDTO)
}

func (h *userHandlerImpl) UpdateProfileImage(c *gin.Context) {
	claims, exists := utils.GetCurrentUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User claims not found in context. Authentication required."})
		return
	}
	userID := claims.UserID

	file, err := c.FormFile("profile_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get profile image file", "details": err.Error()})
		return
	}

	var profileURL string
	uploadDir := "./uploads/profile_images"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			log.Printf("Failed to create upload directory %s: %v", uploadDir, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
	}

	fileExtension := filepath.Ext(file.Filename)
	newUUID := uuid.New().String()

	cleanedFilename := utils.SanitizeFilename(file.Filename)
	baseFilename := strings.TrimSuffix(cleanedFilename, fileExtension)

	uniqueFileName := fmt.Sprintf("%s_%s%s", newUUID, baseFilename, fileExtension)

	filePath := fmt.Sprintf("%s/%s", uploadDir, uniqueFileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Failed to save uploaded file %s: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image", "details": err.Error()})
		return
	}

	profileURL = fmt.Sprintf("%s/uploads/profile_images/%s", h.baseURL, uniqueFileName)

	log.Printf("User %d uploaded file: %s. Stored at: %s. Public URL: %s", userID, file.Filename, filePath, profileURL)

	err = h.userService.UpdateProfileURL(c.Request.Context(), userID, profileURL)
	if err != nil {
		log.Printf("Error updating profile URL for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile image URL", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile image updated successfully", "profile_url": profileURL})
}

func (h *userHandlerImpl) GetStudentSummary(c *gin.Context) {
	summary, err := h.userService.GetOverallStudentSummary(c.Request.Context())
	if err != nil {
		log.Printf("Error getting student summary from service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *userHandlerImpl) CreateAdmin(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or validation failed", "details": err.Error()})
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Class:    req.Class,
		Birthday: req.Birthday,
	}

	err := h.userService.CreateAdmin(c.Request.Context(), user)
	if err != nil {
		switch err.Error() {
		case fmt.Sprintf("user with email %s not found", user.Email):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case fmt.Sprintf("user with email %s already exists", user.Email):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case "email cannot be empty":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "password cannot be empty for hashing":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Printf("Error creating admin user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin user created successfully"})
}

func (h *userHandlerImpl) GetAdminSummary(c *gin.Context) {
	summary, err := h.userService.GetAdminSummary(c.Request.Context())
	if err != nil {
		log.Printf("Error getting admin summary from service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get admin summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *userHandlerImpl) GetMahasiswaUsers(c *gin.Context) {

	mahasiswaUsers, err := h.userService.GetMahasiswaUsers(c.Request.Context())
	if err != nil {
		log.Printf("Error getting mahasiswa users from service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get mahasiswa users"})
		return
	}

	c.JSON(http.StatusOK, mahasiswaUsers)
}

func (h *userHandlerImpl) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format", "details": err.Error()})
		return
	}

	ctx := c.Request.Context()

	err = h.userService.DeleteUser(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User with ID %d not found", userID)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User with ID %d deleted successfully", userID)})
}

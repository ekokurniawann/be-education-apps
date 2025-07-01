package handler

import (
	"be-education/dto"
	"be-education/models"
	"be-education/service"
	"be-education/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userChapterHandlerImpl struct {
	userChapterService service.UserChapterService
}

func NewUserChapterHandler(userChapterService service.UserChapterService) *userChapterHandlerImpl {
	return &userChapterHandlerImpl{userChapterService: userChapterService}
}

func (h *userChapterHandlerImpl) CreateUserChapter(c *gin.Context) {
	var req dto.CreateUserChapterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or validation failed", "details": err.Error()})
		return
	}

	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated or claims not found"})
		return
	}

	userChapter := &models.UserChapter{
		UserID:      claims.UserID,
		ChapterID:   req.ChapterID,
		CompletedAt: req.CompletedAt,
		QuizScore:   req.QuizScore,
	}

	err := h.userChapterService.CreateUserChapter(c.Request.Context(), userChapter)
	if err != nil {
		switch err.Error() {
		case "user ID cannot be zero":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authenticated user ID is invalid"})
		case "chapter ID cannot be zero":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user chapter", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User chapter created successfully", "id": userChapter.ID})
}

func (h *userChapterHandlerImpl) GetUserQuizScores(c *gin.Context) {
	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated or claims not found"})
		return
	}

	quizScores, err := h.userChapterService.GetUserQuizScoresByUserID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user quiz scores", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User quiz scores retrieved successfully", "data": quizScores})
}

func (h *userChapterHandlerImpl) CheckUserChapterCompletion(c *gin.Context) {
	var req dto.CheckChapterCompletionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or validation failed", "details": err.Error()})
		return
	}

	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated or claims not found"})
		return
	}
	userID := claims.UserID
	chapterID := req.ChapterID

	if chapterID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chapter ID must be a positive integer"})
		return
	}

	ctx := c.Request.Context()
	isCompleted, err := h.userChapterService.CheckUserChapterCompleted(ctx, userID, chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check chapter completion status", "details": err.Error()})
		return
	}

	message := "Chapter completion status retrieved successfully"
	if isCompleted {
		message = fmt.Sprintf("User has completed chapter %d", chapterID)
	} else {
		message = fmt.Sprintf("User has not completed chapter %d yet", chapterID)
	}

	c.JSON(http.StatusOK, gin.H{
		"is_completed": isCompleted,
		"chapter_id":   chapterID,
		"user_id":      userID,
		"message":      message,
	})
}

func (h *userChapterHandlerImpl) GetAllUsersChapterScores(c *gin.Context) {

	summary, err := h.userChapterService.GetAllUsersChapterScoresSummary(c.Request.Context())
	if err != nil {
		log.Printf("Error getting all users chapter scores summary: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user chapter scores summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

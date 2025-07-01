package service

import (
	"be-education/dto"
	"be-education/models"
	user_repository "be-education/repository"
	"be-education/utils"
	"context"
	"fmt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	UpdateProfileURL(ctx context.Context, userID int64, profileURL string) error
	GetOverallStudentSummary(ctx context.Context) (*dto.StudentSummary, error)
	CreateAdmin(ctx context.Context, user *models.User) error
	GetMahasiswaUsers(ctx context.Context) ([]*dto.UserResponse, error)
	GetAdminSummary(ctx context.Context) (*dto.AdminSummary, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userServiceImpl struct {
	userRepo user_repository.UserRepository
	jwtUtil  *utils.JWTUtil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	// _, err := s.userRepo.GetUserByID(ctx, id)
	// if err != nil {
	//     return fmt.Errorf("user with ID %d not found: %w", id, err)
	// }

	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("service failed to delete user with ID %d: %w", id, err)
	}
	return nil
}

func NewUserService(userRepo user_repository.UserRepository, jwtUtil *utils.JWTUtil) UserService {
	return &userServiceImpl{userRepo: userRepo, jwtUtil: jwtUtil}
}

func (s *userServiceImpl) CreateAdmin(ctx context.Context, user *models.User) error {
	user.Role = "admin"
	return s.CreateUser(ctx, user)
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *models.User) error {
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	existingUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil && err.Error() != fmt.Sprintf("user with email %s not found", user.Email) {
		return fmt.Errorf("failed to check for existing user: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("service failed to create user: %w", err)
	}
	return nil
}

func (s *userServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with email %s not found", email) {
			return "", fmt.Errorf("invalid credentials")
		}
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	tokenString, err := s.jwtUtil.GenerateJWTToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate authentication token: %w", err)
	}

	return tokenString, nil
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with ID %d not found", id) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve user by ID %d: %w", id, err)
	}

	var userClass string
	if user.Class != nil {
		userClass = *user.Class
	}

	var userBirthday string
	if user.Birthday != nil {
		userBirthday = user.Birthday.Format("2006-01-02")
	}

	var userProfileURL string
	if user.ProfileURL != nil {
		userProfileURL = *user.ProfileURL
	}

	responseDTO := &dto.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Class:      userClass,
		Birthday:   userBirthday,
		Role:       user.Role,
		ProfileURL: userProfileURL,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	return responseDTO, nil
}

func (s *userServiceImpl) UpdateProfileURL(ctx context.Context, userID int64, profileURL string) error {
	err := s.userRepo.UpdateProfileURL(ctx, userID, profileURL)
	if err != nil {
		return fmt.Errorf("service failed to update profile URL: %w", err)
	}
	return nil
}

func (s *userServiceImpl) GetOverallStudentSummary(ctx context.Context) (*dto.StudentSummary, error) {
	classCounts, err := s.userRepo.GetStudentCountsByClass(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get student counts from repository: %w", err)
	}

	totalStudents := 0
	for _, count := range classCounts {
		totalStudents += count
	}

	summary := &dto.StudentSummary{
		TotalStudents: totalStudents,
		ClassCounts:   classCounts,
	}

	return summary, nil
}

func (s *userServiceImpl) GetAdminSummary(ctx context.Context) (*dto.AdminSummary, error) {
	admins, err := s.userRepo.GetAdmins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin users from repository: %w", err)
	}

	totalAdmins, err := s.userRepo.GetTotalAdmins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total admin count from repository: %w", err)
	}

	adminResponses := make([]dto.UserResponse, len(admins))
	for i, admin := range admins {
		var adminClass string
		if admin.Class != nil {
			adminClass = *admin.Class
		}
		var adminBirthday string
		if admin.Birthday != nil {
			adminBirthday = admin.Birthday.Format("2006-01-02")
		}
		var adminProfileURL string
		if admin.ProfileURL != nil {
			adminProfileURL = *admin.ProfileURL
		}

		adminResponses[i] = dto.UserResponse{
			ID:         admin.ID,
			Name:       admin.Name,
			Email:      admin.Email,
			Class:      adminClass,
			Birthday:   adminBirthday,
			Role:       admin.Role,
			ProfileURL: adminProfileURL,
			CreatedAt:  admin.CreatedAt,
			UpdatedAt:  admin.UpdatedAt,
		}
	}

	summary := &dto.AdminSummary{
		TotalAdmins: totalAdmins,
		Admins:      adminResponses,
	}

	return summary, nil
}

func (s *userServiceImpl) GetMahasiswaUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	mahasiswaUsers, err := s.userRepo.GetMahasiswaUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get mahasiswa users from repository: %w", err)
	}

	mahasiswaResponses := make([]*dto.UserResponse, len(mahasiswaUsers))
	for i, user := range mahasiswaUsers {
		var userClass string
		if user.Class != nil {
			userClass = *user.Class
		}
		var userBirthday string
		if user.Birthday != nil {
			userBirthday = user.Birthday.Format("2006-01-02")
		}
		var userProfileURL string
		if user.ProfileURL != nil {
			userProfileURL = *user.ProfileURL
		}

		mahasiswaResponses[i] = &dto.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Class:      userClass,
			Birthday:   userBirthday,
			Role:       user.Role,
			ProfileURL: userProfileURL,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		}
	}

	return mahasiswaResponses, nil
}

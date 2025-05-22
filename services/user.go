// services/user.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/khekrn/apprunner-fiber/models"
)

type UserService struct {
	s3Service *S3Service
}

func NewUserService(s3Service *S3Service) *UserService {
	return &UserService{
		s3Service: s3Service,
	}
}

func (us *UserService) getUserKey(userID string) string {
	return fmt.Sprintf("users/%s/profile.json", userID)
}

func (us *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Name:      req.Name,
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := us.saveUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

func (us *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	key := us.getUserKey(userID)

	exists, err := us.s3Service.ObjectExists(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	data, _, err := us.s3Service.GetObject(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &user, nil
}

func (us *UserService) UpdateUser(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	user, err := us.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Metadata != nil {
		user.Metadata = req.Metadata
	}
	user.UpdatedAt = time.Now()

	if err := us.saveUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (us *UserService) DeleteUser(ctx context.Context, userID string) error {
	// First check if user exists
	_, err := us.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	// Delete user profile
	userKey := us.getUserKey(userID)
	if err := us.s3Service.DeleteObject(ctx, userKey); err != nil {
		return fmt.Errorf("failed to delete user profile: %w", err)
	}

	// Delete all user files
	userFilesPrefix := fmt.Sprintf("users/%s/files/", userID)
	files, err := us.s3Service.ListObjects(ctx, userFilesPrefix)
	if err != nil {
		return fmt.Errorf("failed to list user files: %w", err)
	}

	for _, file := range files {
		if err := us.s3Service.DeleteObject(ctx, file.Key); err != nil {
			return fmt.Errorf("failed to delete user file %s: %w", file.Key, err)
		}
	}

	return nil
}

func (us *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	files, err := us.s3Service.ListObjects(ctx, "users/")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var users []*models.User
	for _, file := range files {
		// Only process profile.json files
		if len(file.Key) > 12 && file.Key[len(file.Key)-12:] == "profile.json" {
			data, _, err := us.s3Service.GetObject(ctx, file.Key)
			if err != nil {
				continue // Skip corrupted files
			}

			var user models.User
			if err := json.Unmarshal(data, &user); err != nil {
				continue // Skip corrupted files
			}

			users = append(users, &user)
		}
	}

	return users, nil
}

func (us *UserService) UserExists(ctx context.Context, userID string) (bool, error) {
	key := us.getUserKey(userID)
	return us.s3Service.ObjectExists(ctx, key)
}

func (us *UserService) saveUser(ctx context.Context, user *models.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	key := us.getUserKey(user.ID)
	_, err = us.s3Service.PutObject(ctx, key, data, "application/json", nil)
	return err
}

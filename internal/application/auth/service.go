package auth

import (
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("crontab-go-secret-key-2024") // 在生产环境中应该从环境变量读取

type Service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

// Login 用户登录
func (s *Service) Login(req *entity.LoginRequest) (*entity.LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, errors.New("用户已被禁用")
	}

	// 生成JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	return &entity.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register 用户注册
func (s *Service) Register(req *entity.RegisterRequest) (*entity.User, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     "user",
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	return user, nil
}

// ValidateToken 验证JWT token
func (s *Service) ValidateToken(tokenString string) (*entity.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return nil, errors.New("用户不存在")
		}
		return user, nil
	}

	return nil, errors.New("无效的token")
}

// generateToken 生成JWT token
func (s *Service) generateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
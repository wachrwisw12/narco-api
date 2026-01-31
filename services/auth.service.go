package services

import (
	"context"
	"errors"
	"time"

	"api-naco/config"
	"api-naco/db"
	middlewares "api-naco/midleware"
	"api-naco/models"

	"golang.org/x/crypto/bcrypt"
)

func AuthRegisterService(user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if len(user.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	// 1️⃣ hash password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (username, password, role_id, fullname)
		VALUES ($1, $2, $3, $4)
		RETURNING username, role_id, fullname
	`

	// 2️⃣ insert + returning
	err = db.DB.QueryRow(
		ctx,
		query,
		user.Username,
		hashedPassword,
		user.RoleId,
		user.Fullname,
	).Scan(
		&user.Username,
		&user.RoleId,
		&user.Fullname,
	)
	if err != nil {
		return nil, err
	}

	// 3️⃣ ไม่ส่ง password กลับ
	user.Password = "d"

	return &user, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost, // cost = 10
	)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func AuthLoginService(
	cfg *config.Config,
	req models.AuthRequest,
) (*models.AuthResponse, error) {
	user, err := findUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("ไม่พบผู้ใช้งาน")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("รหัสผ่านไม่ถูกต้อง")
	}

	token, err := middlewares.GenerateJWT(cfg, user.ID, user.Role)
	if err != nil {
		return nil, errors.New("ไม่สามารถสร้าง token ได้")
	}

	user.Password = ""

	return &models.AuthResponse{
		Token: token,
		Role:  user.Role,
		User:  *user,
	}, nil
}

func findUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT u.id, u.username, u.password, u.role_id ,r.role_name
	FROM users u
    LEFT JOIN roles r ON r.id=u.role_id
	WHERE username = $1
	LIMIT 1
	`

	var u models.User
	err := db.DB.QueryRow(ctx, query, username).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.RoleId,
		&u.Role,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

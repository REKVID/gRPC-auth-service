package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex"`
	Password string
}

type Store struct {
	DB *gorm.DB
}

func NewStore(dsn string) (*Store, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{})
	return &Store{DB: db}, nil
}

func (s *Store) CreateUser(email, password string) error {
	user := User{Email: email, Password: password}
	return s.DB.Create(&user).Error
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := s.DB.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

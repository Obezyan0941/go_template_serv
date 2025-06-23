package db_manager

import (
	"context"
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"golang.org/x/crypto/bcrypt"
)

func NewPostgresConnection(db_data DBConfig) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     db_data.Addr,
		User:     db_data.User,
		Password: db_data.Password,
		Database: db_data.Database,
	})

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func CreateSchema(db *pg.DB, model interface{}) error {
	err := db.Model(model).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Table created")
	return nil
}

func TableExists(db *pg.DB, tableName string) (bool, error) {
	var exists bool
	_, err := db.QueryOne(pg.Scan(&exists), `
        SELECT EXISTS (
            SELECT 1 
            FROM pg_tables 
            WHERE schemaname = 'public' 
            AND tablename = ?
        )`, tableName)
	return exists, err
}

func GetUserDataByID(userID int, db *pg.DB) (*User, error) {
	user := &User{Id: int64(userID)}
	err := db.Model(user).WherePK().Select()
	return user, err
}

func GetUserByName(db *pg.DB, name string) (*User, error) {
	user := new(User)
	err := db.Model(user).
		Where("name = ?", name).
		Select()

	return user, err
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) // returns nil if matches
}

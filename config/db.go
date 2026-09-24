package config

import (
	"fmt"
	"log"

	"go-todo-api/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	sslmode := Config.DBSSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		Config.DBHost,
		Config.DBUser,
		Config.DBPassword,
		Config.DBName,
		Config.DBPort,
		sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed connecting to database:", err)
	}

	DB = db

	err = DB.AutoMigrate(
		&model.User{},
		&model.Todo{},
		&model.Note{},
		&model.Transaction{},
		&model.TransactionCategory{},
		&model.Budget{},
		&model.RecurringTransaction{},
		&model.Alert{},
		&model.UserSession{},
		&model.Habit{},
		&model.HabitLog{},
	)
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	log.Println("database connected successfully")
}

func ConnectTestDB() {
	connect(
		Config.TestDBHost,
		Config.TestDBPort,
		Config.TestDBUser,
		Config.TestDBPassword,
		Config.TestDBName,
		Config.TestDBSSLMode,
	)

}

func connect(
	host string,
	port string,
	user string,
	password string,
	dbname string,
	sslmode string,
) {
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host,
		user,
		password,
		dbname,
		port,
		sslmode,
	)

	database, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal("failed connecting to database:", err)
	}

	DB = database
}

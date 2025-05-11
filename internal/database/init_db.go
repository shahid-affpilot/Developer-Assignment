package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func ConnDB(db_env config.DatabaseConfig) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db_env.Host, strconv.Itoa(db_env.Port), db_env.User, db_env.Password, db_env.Name, db_env.SSLmode)

	fmt.Println(db_env)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB ping failed:", err)
	}
	fmt.Println("Connected to the database.")

	sqlFilePath := "/app/migrations/000001_init_schema/up.sql"
	//sqlFilePath := "/home/shahid/Desktop/Developer-Assignment/migrations/000001_init_schema/up.sql"

	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Println("failed to read sqlFilePath")
	}

	_, err = DB.Exec(string(sqlBytes))

	if err != nil {
		log.Println("DB execute error")
		return
	}
	log.Println("Hurray, database created..")
}

func InitAdminUser(admin config.AdminConfig) {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", admin.Email).Scan(&exists)
	if err != nil {
		log.Fatal("DB check failed:", err)
	}

	salt_pass := os.Getenv("PASSWORD_SALT")
	if !exists {
		hashed, err := bcrypt.GenerateFromPassword([]byte(admin.Password+salt_pass), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Hash error:", err)
		}

		_, err = DB.Exec(`INSERT INTO users (username, email, password_hash, user_type) VALUES ($1, $2, $3, $4)`, admin.Username, admin.Email, string(hashed), admin.UserType)

		if err != nil {
			log.Fatal("Admin insert failed:", err)
		}

		fmt.Println("System admin registered.")
	} else {
		fmt.Println("System admin already exists.")
	}
}

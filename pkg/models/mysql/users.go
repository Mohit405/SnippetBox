package mysql

import (
	"database/sql"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mohit405/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name,email,hashed_password,created) VALUES (?,?,?,UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "u") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	var (
		hashedPassword []byte
		id             int
	)
	stmt := `Select id, hashed_password FROM users WHERE email = ?`

	rows := m.DB.QueryRow(stmt, email)

	err := rows.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return 0, err
	}

	return id, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}

	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`
	
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

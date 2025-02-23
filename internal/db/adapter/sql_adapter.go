package adapter

import (
	"context"
	"fmt"
	"pet-store/internal/models"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
	petTable = "pet"
	tagsTable = "tags"
	petTagsTable = "pet_tags"
	ordersTable = "orders"
	categoryTable = "category"
)

// SQLAdapter - адаптер для работы с БД
type SQLAdapter struct {
	db *sqlx.DB
}

// NewSqlAdapter - конструктор адаптера для работы с БД
func NewSqlAdapter(db *sqlx.DB) *SQLAdapter {
	return &SQLAdapter{db: db}
}

// Create - создание пользователя в БД
func (s *SQLAdapter) Create(ctx context.Context, u *models.UserDTO) (int, error) {
	var id int
	// Проверяем, существует ли уже пользователь с таким именем
	queryCheck := fmt.Sprintf("SELECT id FROM %s WHERE username=$1", usersTable)

	row := s.db.QueryRow(queryCheck, u.UserName)
	err := row.Scan(&id)

	// Если пользователь с таким именем уже существует, возвращаем ошибку
	if err == nil {
		return 0, fmt.Errorf("user with username %s already exists", u.UserName.NullString.String)
	}
	// Если пользователь не существует, создаём нового
	queryInsert := fmt.Sprintf("INSERT INTO %s (username, password, phone, email, firstname, lastname) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", usersTable)
	err = s.db.QueryRowContext(ctx, queryInsert, u.UserName, u.Password, u.Phone, u.Email, u.FirstName, u.LastName).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

func (s *SQLAdapter) GetUserByEmail(ctx context.Context, user *models.UserDTO, email string) error {
	// SQL-запрос для поиска пользователя по email
	query := fmt.Sprintf(`SELECT id, username, email, password, firstname, lastname, phone FROM %s WHERE email = $1`, usersTable)

	// Выполнение запроса к базе данных
	row := s.db.QueryRowContext(ctx, query, email)

	// Маппинг результата в структуру UserDTO
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Phone)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	return nil
}

func (s *SQLAdapter) GetUserByUsername(ctx context.Context, user *models.UserDTO, username string) error {
	// SQL-запрос для поиска пользователя по email
	query := fmt.Sprintf(`SELECT id, username, email, password, firstname, lastname, phone FROM %s WHERE username = $1`, usersTable)

	// Выполнение запроса к базе данных
	row := s.db.QueryRowContext(ctx, query, username)

	// Маппинг результата в структуру UserDTO
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Phone)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	return nil
}

func (s *SQLAdapter) UpdateUser(ctx context.Context, u *models.UserDTO) error {
	query := fmt.Sprintf("UPDATE %s SET firstname = $1, lastname = $2, phone = $3, password = $4 WHERE username = $5", usersTable)

	res, err := s.db.ExecContext(ctx, query, u.FirstName, u.LastName, u.Phone, u.Password, u.UserName)
	if err != nil {
		return err
	}

	// Проверяем, были ли затронуты строки
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with the username %s", u.UserName.String)
	}
	return nil
}

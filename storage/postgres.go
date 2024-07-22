package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"book-catalog/models"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(host, port, user, password, dbname string) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Init() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS books (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            author TEXT NOT NULL,
            isbn TEXT NOT NULL
        )
    `)
	return err
}

func (s *PostgresStorage) GetAll() ([]models.Book, error) {
	rows, err := s.db.Query("SELECT id, title, author, isbn FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (s *PostgresStorage) Get(id string) (models.Book, error) {
	var b models.Book
	err := s.db.QueryRow("SELECT id, title, author, isbn FROM books WHERE id = $1", id).
		Scan(&b.ID, &b.Title, &b.Author, &b.ISBN)
	if err == sql.ErrNoRows {
		return b, fmt.Errorf("book not found")
	}
	return b, err
}

func (s *PostgresStorage) Create(book models.Book) error {
	_, err := s.db.Exec("INSERT INTO books (id, title, author, isbn) VALUES ($1, $2, $3, $4)",
		book.ID, book.Title, book.Author, book.ISBN)
	return err
}

func (s *PostgresStorage) Update(book models.Book) error {
	_, err := s.db.Exec("UPDATE books SET title = $2, author = $3, isbn = $4 WHERE id = $1",
		book.ID, book.Title, book.Author, book.ISBN)
	return err
}

func (s *PostgresStorage) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM books WHERE id = $1", id)
	return err
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

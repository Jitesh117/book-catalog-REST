package storage

import (
	"errors"
	"sync"

	"book-catalog/models"
)

type MemoryStorage struct {
	books map[string]models.Book
	mutex sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		books: make(map[string]models.Book),
	}
}

func (s *MemoryStorage) GetAll() []models.Book {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	books := make([]models.Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return books
}

func (s *MemoryStorage) Get(id string) (models.Book, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	book, ok := s.books[id]
	if !ok {
		return models.Book{}, errors.New("book not found")
	}
	return book, nil
}

func (s *MemoryStorage) Create(book models.Book) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.books[book.ID]; exists {
		return errors.New("book already exists")
	}

	s.books[book.ID] = book
	return nil
}

func (s *MemoryStorage) Update(book models.Book) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.books[book.ID]; !exists {
		return errors.New("book not found")
	}
	s.books[book.ID] = book
	return nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.books[id]; !exists {
		return errors.New("book not found")
	}
	delete(s.books, id)
	return nil
}

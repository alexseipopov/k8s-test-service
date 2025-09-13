package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Record представляет запись для вставки в БД
type Record struct {
	ID        int       `db:"id" json:"id"`
	Message   string    `db:"message" json:"message"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}

// DatabaseManager управляет подключением к БД
type DatabaseManager struct {
	db *sqlx.DB
}

// NewDatabaseManager создает новый экземпляр DatabaseManager
func NewDatabaseManager(connectionString string) (*DatabaseManager, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping БД: %v", err)
	}

	dm := &DatabaseManager{db: db}

	// Создаем таблицу если она не существует
	if err := dm.createTable(); err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	return dm, nil
}

// createTable создает таблицу для хранения записей
func (dm *DatabaseManager) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS records (
		id SERIAL PRIMARY KEY,
		message TEXT NOT NULL,
		timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err := dm.db.Exec(query)
	return err
}

// InsertRecord вставляет новую запись в БД
func (dm *DatabaseManager) InsertRecord(message string) error {
	query := `INSERT INTO records (message, timestamp) VALUES ($1, $2)`
	_, err := dm.db.Exec(query, message, time.Now())
	return err
}

// GetRecordsCount возвращает количество записей в БД
func (dm *DatabaseManager) GetRecordsCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM records`
	err := dm.db.Get(&count, query)
	return count, err
}

// GetRecentRecords возвращает последние N записей
func (dm *DatabaseManager) GetRecentRecords(limit int) ([]Record, error) {
	var records []Record
	query := `SELECT id, message, timestamp FROM records ORDER BY timestamp DESC LIMIT $1`
	err := dm.db.Select(&records, query, limit)
	return records, err
}

// Close закрывает подключение к БД
func (dm *DatabaseManager) Close() error {
	return dm.db.Close()
}

func main() {
	// Получаем строку подключения из переменной окружения или используем значения по умолчанию
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		connectionString = "host=localhost port=5432 user=ticker_user password=ticker_password dbname=ticker_db sslmode=disable"
	}

	// Создаем менеджер БД
	dbManager, err := NewDatabaseManager(connectionString)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer dbManager.Close()

	fmt.Println("Программа запущена. Записи будут добавляться каждые 5 секунд...")
	fmt.Printf("Подключение к БД: %s\n", connectionString)
	fmt.Println("Нажмите Ctrl+C для остановки")

	// Создаем тикер для отправки записей каждые 5 секунд
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Счетчик для уникальных сообщений
	counter := 1

	// Основной цикл
	for {
		select {
		case <-ticker.C:
			message := fmt.Sprintf("Запись #%d - %s", counter, time.Now().Format("2006-01-02 15:04:05"))
			log.Println(message)
			// Вставляем запись в БД
			if err := dbManager.InsertRecord(message); err != nil {
				log.Printf("Ошибка вставки записи: %v", err)
				continue
			}

			// Получаем общее количество записей
			totalCount, err := dbManager.GetRecordsCount()
			if err != nil {
				log.Printf("Ошибка получения количества записей: %v", err)
			} else {
				fmt.Printf("✓ Добавлена запись: %s (всего записей: %d)\n", message, totalCount)
			}

			// Показываем последние 3 записи каждые 5 записей
			if counter%5 == 0 {
				recentRecords, err := dbManager.GetRecentRecords(3)
				if err != nil {
					log.Printf("Ошибка получения последних записей: %v", err)
				} else {
					fmt.Println("Последние записи:")
					for _, record := range recentRecords {
						fmt.Printf("  - ID: %d, Сообщение: %s, Время: %s\n",
							record.ID, record.Message, record.Timestamp.Format("2006-01-02 15:04:05"))
					}
					fmt.Println()
				}
			}

			counter++
		}
	}
}

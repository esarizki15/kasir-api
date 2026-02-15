package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := []models.TransactionDetail{}

	for _, item := range items {
		var name string
		var price, stock int

		// 🔒 lock row
		err := tx.QueryRow(`
			SELECT name, price, stock
			FROM products
			WHERE id = $1
			FOR UPDATE
		`, item.ProductID).Scan(&name, &price, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for %s", name)
		}

		subtotal := price * item.Quantity
		totalAmount += subtotal

		// update stock
		_, err = tx.Exec(`
			UPDATE products
			SET stock = stock - $1
			WHERE id = $2
		`, item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert transaction + ambil created_at
	var transactionID int
	var createdAt time.Time

	err = tx.QueryRow(`
		INSERT INTO transactions (total_amount)
		VALUES ($1)
		RETURNING id, created_at
	`, totalAmount).Scan(&transactionID, &createdAt)

	if err != nil {
		return nil, err
	}

	// insert details + ambil id masing-masing
	for i := range details {
		err = tx.QueryRow(`
			INSERT INTO transaction_details
			(transaction_id, product_id, quantity, subtotal)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).
			Scan(&details[i].ID)

		if err != nil {
			return nil, err
		}

		details[i].TransactionID = transactionID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt,
		Details:     details,
	}, nil
}

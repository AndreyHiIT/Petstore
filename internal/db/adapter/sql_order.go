package adapter

import (
	"context"
	"fmt"
	"pet-store/internal/models"
	"time"
)

func (s *SQLAdapter) CreateOrder(ctx context.Context, order models.Order) (int, error) {

	query := fmt.Sprintf(`
	INSERT INTO %s (petid, quantity, shipdate, status, complete) 
	VALUES($1, $2, $3, $4, $5) 
	RETURNING id`, ordersTable)

	var orderID int
	err := s.db.QueryRowContext(ctx, query, order.PetID, order.Quantity, order.ShipDate, order.Status, order.Complete).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("error createOrder, %v", err)
	}

	return orderID, nil
}

func (s *SQLAdapter) FindOrderByID(ctx context.Context, orderID int) (models.Order, error) {
	query := fmt.Sprintf(`
	SELECT id, petid, quantity, shipdate, status, complete 
	FROM %s 
	WHERE id = $1`, ordersTable)

	var id, petid, quantity int
	var status string
	var shipdate time.Time
	var complete bool

	err := s.db.QueryRowContext(ctx, query, orderID).Scan(&id, &petid, &quantity, &shipdate, &status, &complete)
	if err != nil {
		return models.Order{}, err
	}

	return models.Order{
		ID:       id,
		PetID:    petid,
		Quantity: quantity,
		ShipDate: shipdate,
		Status:   status,
		Complete: complete,
	}, nil
}

func (s *SQLAdapter) DeleteOrderByID(ctx context.Context, orderID int) error {
	query := fmt.Sprintf(`
	DELETE FROM %s 
	WHERE id = $1`, ordersTable)

	result, err := s.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return err
	}

	// Проверяем, были ли затронуты строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no order found with ID")
	}

	return nil
}

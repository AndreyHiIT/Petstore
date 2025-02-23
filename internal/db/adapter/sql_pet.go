package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	myerrors "pet-store/internal/infrastructure/errors"
	"pet-store/internal/models"
	"strings"
)

func (s *SQLAdapter) UpdatePetForm(ctx context.Context, name, status string, petID int) error {
	query := fmt.Sprintf("UPDATE %s Set", petTable)
	var params []interface{}
	var setClauses []string
	count := 1

	if name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", count))
		count++
		params = append(params, name)
	}
	if status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", count))
		count++
		params = append(params, status)
	}

	query += " " + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", count)
	params = append(params, petID)

	result, err := s.db.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Если строки не были затронуты, возвращаем ошибку
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated, possibly invalid ID")
	}

	return nil
}

func (s *SQLAdapter) FindPetbyID(ctx context.Context, petid int) (models.Pet, error) {
	query := `
		SELECT 
			p.id AS pet_id,
			p.name AS pet_name,
			p.photourls,
			p.status AS pet_status,
			c.id,
			c.name AS category_name,
			t.id,
			t.name AS tag_name
		FROM pet p
		JOIN category c ON p.category = c.id
		LEFT JOIN pet_tags pt ON p.id = pt.pet_id
		LEFT JOIN tags t ON pt.tag_id = t.id
		WHERE p.id IN ($1)
		ORDER BY p.id;
	`

	rows, err := s.db.QueryContext(ctx, query, petid)
	if err != nil {
		return models.Pet{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return models.Pet{}, myerrors.ErrPetNotFound
	}
	var currentPet models.Pet

	for rows.Next() {
		var petID, categoryID, tagID int
		var petName, photourls, petStatus, categoryName, tagName string

		// Сканирование строки
		if err := rows.Scan(&petID, &petName, &photourls, &petStatus, &categoryID, &categoryName, &tagID, &tagName); err != nil {
			return models.Pet{}, fmt.Errorf("failed to scan row: %w", err)
		}
		// Добавляем питомца в переменную currentPet
		if currentPet.ID != petID {
			currentPet = models.Pet{
				ID:        petID,
				Name:      petName,
				PhotoUrls: []string{photourls},
				Status:    petStatus,
				Category: models.Category{
					ID:   categoryID,
					Name: categoryName,
				},
				Tags: []models.Tag{
					{
						ID:   tagID,
						Name: tagName,
					},
				},
			}
		} else {
			// Если есть ещё строки добавляем новые теги
			tag := models.Tag{
				ID:   tagID,
				Name: tagName,
			}
			currentPet.Tags = append(currentPet.Tags, tag)
		}
	}
	return currentPet, nil
}

// Проверяет наличие строки str в таблице tablename в поле colimnname и возвращает id строки если она есть
func (s *SQLAdapter) CheckFields(tablename, str, columnname string) (int, error) {
	query := fmt.Sprintf(
		`SELECT id
		FROM %s
		WHERE %s = $1 LIMIT 1;`, tablename, columnname)

	var id int
	err := s.db.QueryRow(query, str).Scan(&id)
	if err != nil {
		// Если ошибка связана с тем, что строка не найдена
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

func (s *SQLAdapter) AddPet(ctx context.Context, pet *models.Pet) error {
	//Проверяем корректность введённой категории
	idCategory, err := s.CheckFields(categoryTable, pet.Category.Name, "name")
	if err != nil {
		return err
	}
	if idCategory == 0 {
		return errors.New("некорректно введена категория")
	}
	//Вставляем pet в таблицу
	var petID int
	queryPet := fmt.Sprintf(
		`INSERT INTO %s (category, name, photoUrls, status)
				VALUES ($1, $2, $3, $4)
				RETURNING id`, petTable)
	urls := strings.Join(pet.PhotoUrls, ", ")
	err = s.db.QueryRowContext(ctx, queryPet, idCategory, pet.Name, urls, pet.Status).Scan(&petID)
	if err != nil {
		return fmt.Errorf("addpet qurePet, %v", err)
	}

	//Вставим теги в таблицу
	for i := range pet.Tags {
		queryTag := fmt.Sprintf(
			`INSERT INTO %s (name)
			VALUES ($1)
			ON CONFLICT (name) DO NOTHING
			RETURNING id`, tagsTable)
		var tagID int
		err := s.db.QueryRowContext(ctx, queryTag, pet.Tags[i].Name).Scan(&tagID)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return fmt.Errorf("addpet queryTag, %v", err)
		}
		// Если тег уже существует, возвращаем его ID
		if tagID == 0 {
			queryTag = fmt.Sprintf(
				`SELECT id FROM %s WHERE name = $1`, tagsTable)
			err := s.db.QueryRowContext(ctx, queryTag, pet.Tags[i].Name).Scan(&tagID)
			if err != nil {
				return fmt.Errorf("addpet queryTag when taID=0, %v", err)
			}
		}
		//Связываем теги с вставленным питомцем
		queryTagPet := fmt.Sprintf(
			`INSERT INTO %s (pet_id, tag_id)
			VALUES($1, $2)`, petTagsTable)
		_, err = s.db.ExecContext(ctx, queryTagPet, petID, tagID)
		if err != nil {
			return fmt.Errorf("addpet queryTagPet, %v", err)
		}
	}
	return nil
}

func (s *SQLAdapter) UpdatePet(ctx context.Context, pet models.Pet) error {
	//Проверяем корректность введённой категории
	var idCategory int
	if pet.Category.Name != "" {
		idCategory, err := s.CheckFields(categoryTable, pet.Category.Name, "name")
		if err != nil {
			return err
		}
		if idCategory == 0 {
			return errors.New("некорректно введена категория")
		}
	}
	query := fmt.Sprintf(`
	UPDATE %s
	SET
    name = COALESCE(NULLIF($1, ''), name),
    category = COALESCE(NULLIF($2, 0), category),
    status = COALESCE(NULLIF($3, ''), status),
    photourls = COALESCE($4, photourls)
WHERE id = $5;`, petTable)
	urls := strings.Join(pet.PhotoUrls, ", ")
	_, err := s.db.ExecContext(ctx, query, pet.Name, idCategory, pet.Status, urls, pet.ID)
	if err != nil {
		return fmt.Errorf("error in sqlAdapter-(UpdatePet)-execintable: %v", err)
	}
	fmt.Println(pet.Tags)
	if len(pet.Tags) != 0 {
		queryDel := fmt.Sprintf(`
			DELETE FROM %s WHERE pet_id = $1`, petTagsTable)
		_, err = s.db.ExecContext(ctx, queryDel, pet.ID)
		if err != nil {
			return fmt.Errorf("updatepet queryDel: %v", err)
		}
		for i := range pet.Tags {
			queryTag := fmt.Sprintf(
				`INSERT INTO %s (name)
				VALUES ($1)
				ON CONFLICT (name) DO NOTHING
				RETURNING id`, tagsTable)
			var tagID int
			err := s.db.QueryRowContext(ctx, queryTag, pet.Tags[i].Name).Scan(&tagID)
			if err != nil && err.Error() != "sql: no rows in result set" {
				return fmt.Errorf("updatepet queryTag, %v", err)
			}
			// Если тег уже существует, возвращаем его ID
			if tagID == 0 {
				queryTag = fmt.Sprintf(
					`SELECT id FROM %s WHERE name = $1`, tagsTable)
				err := s.db.QueryRowContext(ctx, queryTag, pet.Tags[i].Name).Scan(&tagID)
				if err != nil {
					return fmt.Errorf("updatepet queryTag: %v", err)
				}
			}
			//Связываем теги с обновляемым питомцем
			queryTagPet := fmt.Sprintf(
				`INSERT INTO %s (pet_id, tag_id)
		VALUES($1, $2)`, petTagsTable)
			_, err = s.db.ExecContext(ctx, queryTagPet, pet.ID, tagID)
			if err != nil {
				return fmt.Errorf("addpet queryTagPet, %v", err)
			}
		}
	}
	return nil
}

func (a *SQLAdapter) FindPetbyStatus(ctx context.Context, statuses []string) ([]models.Pet, error) {
	placeholders := make([]string, len(statuses))
	for i := range statuses {
		placeholders[i] = fmt.Sprintf("$%d", i+1) // $1, $2 и т.д.
	}
	query := fmt.Sprintf(`
		SELECT 
			p.id AS pet_id,
			p.name AS pet_name,
			p.photourls,
			p.status AS pet_status,
			c.id,
			c.name AS category_name,
			t.id,
			t.name AS tag_name
		FROM pet p
		JOIN category c ON p.category = c.id
		LEFT JOIN pet_tags pt ON p.id = pt.pet_id
		LEFT JOIN tags t ON pt.tag_id = t.id
		WHERE p.status IN (%s)
		ORDER BY p.id;
	`, strings.Join(placeholders, ", "))

	// Подготовка аргументов запроса
	args := make([]interface{}, len(statuses))
	for i, status := range statuses {
		args[i] = status
	}

	// Выполнение запроса
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var pets []models.Pet
	var currentPet *models.Pet

	for rows.Next() {
		var petID, categoryID, tagID int
		var petName, photourls, petStatus, categoryName, tagName string

		// Сканирование строки
		if err := rows.Scan(&petID, &petName, &photourls, &petStatus, &categoryID, &categoryName, &tagID, &tagName); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Если это новый питомец, добавляем его в список
		if currentPet == nil || currentPet.ID != petID {
			if currentPet != nil {
				pets = append(pets, *currentPet)
			}
			currentPet = &models.Pet{
				ID:        petID,
				Name:      petName,
				PhotoUrls: []string{photourls},
				Status:    petStatus,
				Category: models.Category{
					ID:   categoryID,
					Name: categoryName,
				},
				Tags: []models.Tag{
					{
						ID:   tagID,
						Name: tagName,
					},
				},
			}
		} else {
			// Если это тот же питомец, добавляем новый тег
			tag := models.Tag{
				ID:   tagID,
				Name: tagName,
			}
			currentPet.Tags = append(currentPet.Tags, tag)
		}
	}

	// Добавляем последнего питомца в список
	if currentPet != nil {
		pets = append(pets, *currentPet)
	}

	// Проверка на ошибки при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return pets, nil
}

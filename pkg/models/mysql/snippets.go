package mysql

import (
	"database/sql"

	"github.com/mohit405/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connectino pool.
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title,content,created,expires) 
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	var snip models.Snippet

	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return &snip, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	data:=[]*models.Snippet{}

	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP()
	ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		snipData:=&models.Snippet{}
		err = rows.Scan(&snipData.ID, &snipData.Title, &snipData.Content, &snipData.Created, &snipData.Expires)
		if err != nil {
			return nil, err
		}

		data = append(data, snipData)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

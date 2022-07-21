package mysql

import (
	"database/sql"
	"github.com/Onolax/SnipB/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

//will insert a new snippet
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	//executes the command
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

//will return the specific snippet
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	//command to run sql
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?`
	//executing the command and returns a single row
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	//copying the row to s
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

// will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	//sql statement
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	//executing the stmt and returns multiple rows
	/*basic difference between Query,QueryRow,Exec is that first returns multiple rows
	second return a single row and third return th value which is not a row*/
	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	//close the file before end of Latest()
	defer rows.Close()

	snippets := []*models.Snippet{}

	//iterating over all the rows
	for rows.Next() {
		//pointing to start
		s := &models.Snippet{}
		//copying the rows to the new snippets we created
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}
		//Appended it to the slice of snippets.
		snippets = append(snippets, s)
	}
	//check if any error encountered during the loop
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

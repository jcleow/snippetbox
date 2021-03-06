package mysql

import (
	"database/sql"

	"github.com/jcleow/snippetbox/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool

type SnippetModel struct {
	DB *sql.DB
}

// func (m *SnippetModel)Get(id int)(*models.Snippet, error){

// 	stmt := `SELECT id, title, content, created, expires FROM snippets
// 	WHERE expires > UTC_TIMESTAMP() AND id = ?`

// 	// Use QueryRow() method on connection pool to execute our SQL statement
// 	// passing in the untrusted id variable as the value for the placeholder parameter
// 	// This returns a pointer to a sql.Row object which holds
// 	// the result from the database
// 	row := m.DB.QueryRow(stmt, id)

// 	// Initialize a pointer to a new zeroed Snippet struct
// 	s := &models.Snippet{}

// 	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
// 	if err == sql.ErrNoRows{
// 		return nil, models.ErrNoRecord
// 	} else if err != nil{
// 		return nil, err
// 	}
// 	//if everything went ok, then return the Snippet object.
// 	return s, nil
// }

func (m *SnippetModel) Get(id int)(*models.Snippet, error){
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// Initialize a pointer to a new zeroed Snippet struct
	s := &models.Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

//This will insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error){
	stmt := `INSERT INTO snippets(title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil{
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type before returning
	return int(id) ,nil
}


// This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error){
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	
	// Use the Query() method on the connection pool to execute our 
	// SQL Statement. This returns a sql.Rows result set containing the result of 
	// our query
	rows, err := m.DB.Query(stmt)
	if err != nil{
		return nil, err
	}

	// We defer rows.Close() to ensure that sql.Rows resultset is 
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query() method
	// Otherwise, if Query() returns an error, you'll get a panic trying to close a nil resultset
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects
	snippets := []*models.Snippet{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first ( and then each subsequent) row to be acted on by the 
	// rows.Scan() method. If iteration over all the rows completes then the 
	// resultset automatically closes itself and frees-up the underlying database connection
	for rows.Next(){
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}

		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan
		// must be pointers to the place you want to copy the data into
		// and the number of arguments must be exactly the same as the number of 
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		
		if err != nil{
			return nil, err
		}
		// Append it to the slice of snippets
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve an error
	// that was encountered during the iteration. It's important to call this - don't assume
	// that a successful iteration was completed over the whole resultset.
	if err = rows.Err(); err != nil{
		return nil, err
	}
	// if everything went ok, return the Snippets slice.
	return snippets, nil
}
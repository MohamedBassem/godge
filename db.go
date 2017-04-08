package godge

func (s *Server) initDB() error {
	const schema = `
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username varchar(255),
    password varchar(255)
	);

	CREATE TABLE IF NOT EXISTS scoreboard (
		id INTEGER PRIMARY KEY,
		username INTEGER,
		task_name varchar(255),
		verdict varchar(255),
		submitted_at DATETIME
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

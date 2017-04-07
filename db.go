package godge

func (s *Server) initDB() error {
	const schema = `
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username varchar(255),
    password varchar(255)
	);

	CREATE TABLE IF NOT EXISTS submission (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		task_name varchar(255),
		language varchar(255),
		verdict varchar(255)
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

package database

type TaskInformation struct {
	ID          int
	Name        string
	Description string
	IsFinished  bool
	IsArchived  bool
	IDGroup     int
}

// TaskInit initialize the task table
func (db *DBHandler) TaskInit() error {
	createTaskTable := `CREATE TABLE IF NOT EXISTS stormtask_task (
		id_task INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		is_finished BOOLEAN NOT NULL,
		is_archived BOOLEAN NOT NULL,
		id_group INT REFERENCES stormtask_group(id_group)
	)`
	_, err := db.Handler.Exec(createTaskTable)
	return err
}

// GetTaskByID get the task using his id
// If the task is found, this returns a TaskInformation pointer
// If an error occurred, this returns nil and the error
func (db *DBHandler) GetTaskByID(id int) (*TaskInformation, error) {
	getTaskRequest := "SELECT * FROM stormtask_task WHERE id_task = ?"
	statement, err := db.Handler.Prepare(getTaskRequest)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	row := statement.QueryRow(id)
	if row == nil {
		return nil, nil
	}
	var task TaskInformation
	err = row.Scan(&task.ID, &task.Name, &task.Description, &task.IsFinished, &task.IsArchived, &task.IDGroup)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetTaskByGroup get all the tasks of a selected group
// In the nominal case this return a TaskInformation table pointer
// If an error occurred, this returns nil and the error
// If the group is not found, this returns nil and nil
func (db *DBHandler) GetTaskByGroup(id int) (*[]TaskInformation, error) {
	getTasksByGroup := "SELECT * FROM stormtask_task WHERE id_group = ?"
	statement, err := db.Handler.Prepare(getTasksByGroup)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	rows, err := statement.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []TaskInformation
	for rows.Next() {
		var task TaskInformation
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.IsFinished, &task.IsArchived, &task.IDGroup)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return &tasks, nil
}

// AddTask add a task to the database
// In the nominal case, this returns the task created
// In case of error, this returns nil and the error
func (db *DBHandler) AddTask(
	name,
	description string,
	isFinished,
	isArchived bool,
	group int) (*TaskInformation, error) {
	insertTask := `INSERT INTO stormtask_task(name, description, is_finished, is_archived, id_group) VALUES (?,?,?,?,?)`
	statement, err := db.Handler.Prepare(insertTask)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Exec(name, description, isFinished, isArchived, group)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return db.GetTaskByID(int(id))
}

// ModifyTask modify a task selected by it id
// In the nominal case, this returns a TaskInformation
// In case of error, this returns nil and the error
func (db *DBHandler) ModifyTask(id int,
	name,
	description string,
	isFinished,
	isArchived bool,
	group int) (*TaskInformation, error) {
	modifyTask := "UPDATE stormtask_task SET name = ?, description = ?, is_finished = ?," +
		"is_archived = ?, id_group = ? WHERE id_task = ?"
	statement, err := db.Handler.Prepare(modifyTask)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	_, err = statement.Exec(name, description, isFinished, isArchived, group, id)
	if err != nil {
		return nil, err
	}
	return db.GetTaskByID(id)
}

// DeleteTask delete a task from the database
// If an error occurred, this returns an error, else, it returns nil
func (db *DBHandler) DeleteTask(id int) error {
	deleteTask := "DELETE FROM stormtask_task WHERE id_task = ?"
	statement, err := db.Handler.Prepare(deleteTask)
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	return err
}

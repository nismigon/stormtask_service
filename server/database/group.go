package database

type GroupInformation struct {
	ID     int
	Name   string
	UserID int
}

// GroupInit creates the table containing the group if it doesn't already exists
// If an error occurred, this method returns an error
func (db *DBHandler) GroupInit() error {
	createGroupTable := `CREATE TABLE IF NOT EXISTS stormtask_group (
		id_group INT AUTO_INCREMENT,
		id_user INT,
		name VARCHAR(255),
		PRIMARY KEY (id_group),
		FOREIGN KEY (id_user) REFERENCES stormtask_user (id_user)
	)`
	_, err := db.Handler.Exec(createGroupTable)
	return err
}

// GetGroupByID get a group which correspond to the group id
// Return a group or an error if a database error occurred
func (db *DBHandler) GetGroupByID(id int) (*GroupInformation, error) {
	getGroupRequest := `SELECT id_group, name, id_user FROM stormtask_group WHERE id_group=?`
	statement, err := db.Handler.Prepare(getGroupRequest)
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var group GroupInformation
		err = rows.Scan(&group.ID, &group.Name, &group.UserID)
		if err != nil {
			return nil, err
		}
		return &group, nil
	}
	return nil, nil
}

// GetGroupByUserAndName get a group which correspond to the user and the name
// Return a group or an error if a database error occurred
func (db *DBHandler) GetGroupByUserAndName(user_id int, name string) (*GroupInformation, error) {
	getGroupRequest := `SELECT id_group, name, id_user FROM stormtask_group WHERE id_user=? AND name=?`
	statement, err := db.Handler.Prepare(getGroupRequest)
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(user_id, name)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var group GroupInformation
		err = rows.Scan(&group.ID, &group.Name, &group.UserID)
		if err != nil {
			return nil, err
		}
		return &group, nil
	}
	return nil, nil
}

// AddGroup add a group to the database
// If the group have been added, this returns a GroupInformation
// If an error occurred, this returns an errror
func (db *DBHandler) AddGroup(user_id int, group_name string) (*GroupInformation, error) {
	addUserRequest := `INSERT INTO stormtask_group (id_user, name) VALUES (?, ?)`
	statement, err := db.Handler.Prepare(addUserRequest)
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(user_id, group_name)
	if err != nil {
		return nil, err
	}
	return db.GetGroupByUserAndName(user_id, group_name)
}

// ChangeGroupName change the name of the group
// Return a group object or an error if an error occurred in the database
func (db *DBHandler) ChangeGroupName(group_id int, group_name string) (*GroupInformation, error) {
	changeGroupNameRequest := `UPDATE stormtask_group SET name=? WHERE id_group=?`
	statement, err := db.Handler.Prepare(changeGroupNameRequest)
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(group_name, group_id)
	if err != nil {
		return nil, err
	}
	return db.GetGroupByID(group_id)
}

// DeleteGroup deletes a group from the table
// Return nil if the group has been deleted or an error if an error occurred
func (db *DBHandler) DeleteGroup(group_id int) error {
	deleteRequest := `DELETE FROM stormtask_group WHERE id_group=?`
	statement, err := db.Handler.Prepare(deleteRequest)
	if err != nil {
		return err
	}
	_, err = statement.Exec(group_id)
	return err
}

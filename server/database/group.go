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
		FOREIGN KEY (id_user) REFERENCES stormtask_user (id_user),
		CONSTRAINT group_name_per_user_unique UNIQUE (id_user, name)
	)`
	_, err := db.Handler.Exec(createGroupTable)
	return err
}

// GetGroupByID get a group which correspond to the group id
// Return a group or an error if a database error occurred
func (db *DBHandler) GetGroupByID(id int) (*GroupInformation, error) {
	getGroupRequest := `SELECT id_group, name, id_user FROM stormtask_group WHERE id_group=?`
	statement, err := db.Handler.Prepare(getGroupRequest)
	defer func() {
		_ = statement.Close()
	}()
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(id)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
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
func (db *DBHandler) GetGroupByUserAndName(userID int, name string) (*GroupInformation, error) {
	getGroupRequest := `SELECT id_group, name, id_user FROM stormtask_group WHERE id_user=? AND name=?`
	statement, err := db.Handler.Prepare(getGroupRequest)
	defer func() {
		_ = statement.Close()
	}()
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(userID, name)
	defer func() {
		_ = rows.Close()
	}()
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

// GetGroupsByUserID get all the groups of a selected user
// In the nominal case, this returns a list of groups
// If an error occurred during the process, this returns an error
func (db *DBHandler) GetGroupsByUserID(userID int) (*[]GroupInformation, error) {
	getGroupsRequest := `SELECT id_group, name, id_user FROM stormtask_group WHERE id_user=?`
	statement, err := db.Handler.Prepare(getGroupsRequest)
	defer func() {
		_ = statement.Close()
	}()
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(userID)
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		return nil, err
	}
	var groups []GroupInformation
	for rows.Next() {
		var group GroupInformation
		err = rows.Scan(&group.ID, &group.Name, &group.UserID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return &groups, nil
}

// AddGroup add a group to the database
// If the group have been added, this returns a GroupInformation
// If an error occurred, this returns an errror
func (db *DBHandler) AddGroup(userID int, groupName string) (*GroupInformation, error) {
	addUserRequest := `INSERT INTO stormtask_group (id_user, name) VALUES (?, ?)`
	statement, err := db.Handler.Prepare(addUserRequest)
	defer func() {
		_ = statement.Close()
	}()
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(userID, groupName)
	if err != nil {
		return nil, err
	}
	return db.GetGroupByUserAndName(userID, groupName)
}

// ModifyGroup change the name of the group
// Return a group object or an error if an error occurred in the database
func (db *DBHandler) ModifyGroup(groupID int, groupName string) (*GroupInformation, error) {
	changeGroupNameRequest := `UPDATE stormtask_group SET name=? WHERE id_group=?`
	statement, err := db.Handler.Prepare(changeGroupNameRequest)
	defer func() {
		_ = statement.Close()
	}()
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(groupName, groupID)
	if err != nil {
		return nil, err
	}
	return db.GetGroupByID(groupID)
}

// DeleteGroup deletes a group from the table
// Return nil if the group has been deleted or an error if an error occurred
func (db *DBHandler) DeleteGroup(groupID int) error {
	err := db.DeleteTasksByGroup(groupID)
	if err != nil {
		return err
	}
	deleteRequest := `DELETE FROM stormtask_group WHERE id_group=?`
	statement, err := db.Handler.Prepare(deleteRequest)
	defer func() {
		_ = statement.Close()
	}()
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupID)
	return err
}

// DeleteGroupsByUser deletes all the groups of a selected user
// In the nominal case, this returns a nil error
// If an error occurred during the request to the database, this returns the error generated
func (db *DBHandler) DeleteGroupsByUser(id int) error {
	getGroups, err := db.GetGroupsByUserID(id)
	if err != nil {
		return err
	}
	for _, group := range *getGroups {
		err = db.DeleteGroup(group.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

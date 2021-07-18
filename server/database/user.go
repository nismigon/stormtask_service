package database

type UserInformation struct {
	ID       int
	Name     string
	Email    string
	Password string
	IsAdmin  bool
}

// UserInit creates the table containing the user if it doesn't already exists
// If an error occurred, this method returns an error
func (db *DBHandler) UserInit() error {
	createUserTable := `CREATE TABLE IF NOT EXISTS stormtask_user (
    	id_user INT PRIMARY KEY AUTO_INCREMENT,
    	email VARCHAR(255) UNIQUE,
		name VARCHAR(255),
		password VARCHAR(255),
		is_admin BOOLEAN
	)`
	_, err := db.Handler.Exec(createUserTable)
	return err
}

// GetUserByEmail returns a pointer to a UserInformation
// If the user is not found, this function returns nil
// If the database return an error, this error is propagated
func (db *DBHandler) GetUserByEmail(email string) (*UserInformation, error) {
	getUser := "SELECT id_user, name, email, is_admin FROM stormtask_user WHERE email=?"
	statement, err := db.Handler.Prepare(getUser)
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(email)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var user UserInformation
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

// GetUserById returns a pointer to a UserInformation
// If the user is not found, this function returns nil
// If the database return an error, this error is propagated
func (db *DBHandler) GetUserById(id int) (*UserInformation, error) {
	getUser := "SELECT id_user, name, email, is_admin FROM stormtask_user WHERE id_user=?"
	statement, err := db.Handler.Prepare(getUser)
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var user UserInformation
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

// Authenticate verify if a corresponding user exist in the database
// Return a pointer of UserInformation if a user is found,
// nil if no user with this email ant this password is found or
// error if an error occurred during the request
func (db *DBHandler) Authenticate(email, password string) (*UserInformation, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	if user.Password == password {
		return user, nil
	} else {
		return nil, nil
	}
}

// AddUser add a user to the database
// Return user if the user have been added or an error if the database return an error
func (db *DBHandler) AddUser(email, name, password string, is_admin bool) (*UserInformation, error) {
	addUserRequest := `INSERT INTO stormtask_user (name, email, password, is_admin) VALUES (?, ?, ?, ?)`
	statement, err := db.Handler.Prepare(addUserRequest)
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(name, email, password, is_admin)
	if err != nil {
		return nil, err
	}
	return db.GetUserByEmail(email)
}

// ModifyUser modify a user to the database
// Return the modified user or an error if the database returns an error
func (db *DBHandler) ModifyUser(id int, email, name, password string) (*UserInformation, error) {
	modifyUserRequest := `UPDATE stormtask_user SET email=?, name=?, password=? WHERE id_user=?`
	statement, err := db.Handler.Prepare(modifyUserRequest)
	if err != nil {
		return nil, err
	}
	_, err = statement.Exec(email, name, password, id)
	if err != nil {
		return nil, err
	}
	return db.GetUserById(id)
}

// DeleteUser delete a user of the table
// Return nil if the user have been deleted or err if an error occured
func (db *DBHandler) DeleteUser(id int) error {
	deleteUserRequest := `DELETE FROM stormtask_user WHERE email=?`
	statement, err := db.Handler.Prepare(deleteUserRequest)
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	return err
}

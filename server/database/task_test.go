package database

import (
	"teissem/stormtask/server/configuration"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TaskTestSuite struct {
	suite.Suite
	Handler *DBHandler
	User    *UserInformation
	Group   *GroupInformation
	Task    *TaskInformation
}

func (suite *TaskTestSuite) SetupTest() {
	conf, err := configuration.Parse("../../configuration.json")
	if err != nil {
		suite.T().Errorf("Failed to parse the configuration file : " + err.Error())
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		suite.T().Errorf("Failed to initialize the connection with the database : " + err.Error())
	}
	suite.Handler = handler
	user, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		suite.T().Errorf("Failed to add the user in the database : " + err.Error())
	}
	suite.User = user
	group, err := handler.AddGroup(user.ID, "TestGroup")
	if err != nil {
		suite.T().Errorf("Failed to add the group in the database : " + err.Error())
	}
	suite.Group = group
	task, err := handler.AddTask("TestTask", "Description of the task", false, false, group.ID)
	if err != nil {
		suite.T().Errorf("Failed to add the task in the database : " + err.Error())
	}
	suite.Task = task
}

func (suite *TaskTestSuite) TearDownTest() {
	err := suite.Handler.DeleteUser(suite.User.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete user : " + err.Error())
	}
	_ = suite.Handler.Close()
}

func (suite *TaskTestSuite) TestAddAndDeleteTask() {
	task, err := suite.Handler.AddTask("Task", "A little description", false, false, suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to add the task : " + err.Error())
	}
	if task.Name != "Task" {
		suite.T().Errorf("Failed to set the name of the task")
	}
	if task.Description != "A little description" {
		suite.T().Errorf("Failed to set the description of the task")
	}
	if task.IsArchived {
		suite.T().Errorf("Failed to set the archived status of the task")
	}
	if task.IsFinished {
		suite.T().Errorf("Failed to set the finished status of the task")
	}
	if task.IDGroup != suite.Group.ID {
		suite.T().Errorf("Failed to set the owner group of the task")
	}
	err = suite.Handler.DeleteTask(task.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete the task : " + err.Error())
	}
}

func (suite *TaskTestSuite) TestAddTaskWrongGroup() {
	_, err := suite.Handler.AddTask("Task", "A little description", false, false, -1)
	if err == nil {
		suite.T().Errorf("Failed to add the task : no error when a wrong group is given")
	}
}

func (suite *TaskTestSuite) TestModifyTaskRight() {
	task, err := suite.Handler.ModifyTask(suite.Task.ID, "New Task", "A bigger description", true, true, suite.Task.IDGroup)
	if err != nil {
		suite.T().Errorf("Failed to modify task : " + err.Error())
	}
	if task.Name != "New Task" {
		suite.T().Errorf("Failed to set the name of the task")
	}
	if task.Description != "A bigger description" {
		suite.T().Errorf("Failed to set the description of the task")
	}
	if !task.IsArchived {
		suite.T().Errorf("Failed to set the archived status of the task")
	}
	if !task.IsFinished {
		suite.T().Errorf("Failed to set the finished status of the task")
	}
	if task.IDGroup != suite.Task.IDGroup {
		suite.T().Errorf("Failed to set the owner group of the task")
	}
}

func (suite *TaskTestSuite) TestModifyTaskWrongID() {
	_, err := suite.Handler.ModifyTask(-1, "Test", "Test", false, false, suite.Task.IDGroup)
	if err == nil {
		suite.T().Errorf("Failed to modify task : give a wrong ID doesn't return an error")
	}
}

func (suite *TaskTestSuite) TestModifyTaskWrongGroupID() {
	_, err := suite.Handler.ModifyTask(suite.Task.ID, "Test", "Test", false, false, -1)
	if err == nil {
		suite.T().Errorf("Failed to modify task : give a wrong group ID doesn't return an error")
	}
}

func (suite *TaskTestSuite) TestGetTaskByIDRight() {
	task, err := suite.Handler.GetTaskByID(suite.Task.ID)
	if err != nil {
		suite.T().Errorf("Failed to get task by its ID : " + err.Error())
	}
	if task.Name != suite.Task.Name {
		suite.T().Errorf("Failed to set the name of the task")
	}
	if task.Description != suite.Task.Description {
		suite.T().Errorf("Failed to set the description of the task")
	}
	if task.IsArchived != suite.Task.IsArchived {
		suite.T().Errorf("Failed to set the archived status of the task")
	}
	if task.IsFinished != suite.Task.IsFinished {
		suite.T().Errorf("Failed to set the finished status of the task")
	}
	if task.IDGroup != suite.Task.IDGroup {
		suite.T().Errorf("Failed to set the owner group of the task")
	}
}

func (suite *TaskTestSuite) TestGetTaskByIDWrongID() {
	_, err := suite.Handler.GetTaskByID(-1)
	if err == nil {
		suite.T().Errorf("Failed to get task by its ID : No error whereas wrong task ID given")
	}
}

func (suite *TaskTestSuite) TestGetTaskByGroupRight() {
	_, err := suite.Handler.AddTask("Test", "Description", false, false, suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to add task : " + err.Error())
	}
	tasks, err := suite.Handler.GetTasksByGroup(suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to get tasks by group : " + err.Error())
	}
	if len(*tasks) != 2 {
		suite.T().Errorf("Failed to get tasks by group : The result of the get task by group is different of 1")
	}
}

func (suite *TaskTestSuite) TestDBHandler_DeleteTasksByGroup() {
	_, err := suite.Handler.AddTask("Test1", "Description", false, false, suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to add task : " + err.Error())
	}
	_, err = suite.Handler.AddTask("Test2", "Description", false, false, suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to add task : " + err.Error())
	}
	err = suite.Handler.DeleteTasksByGroup(suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to delete tasks of the group : " + err.Error())
	}
	tasks, err := suite.Handler.GetTasksByGroup(suite.Group.ID)
	if err != nil {
		suite.T().Errorf("Failed to get the tasks for a selected group : " + err.Error())
	}
	if len(*tasks) != 0 {
		suite.T().Errorf("Failed to delete the tasks : The returned array shoul be empty")
	}
}

func (suite *TaskTestSuite) TestUniqueConstraintTask() {
	task, err := suite.Handler.AddTask("TestTask", "An other description", false, false, suite.Group.ID)
	if err == nil {
		suite.T().Errorf("Failed to add task : no error whereas break unique constraint on name")
		err = suite.Handler.DeleteTask(task.ID)
		if err != nil {
			suite.T().Errorf("Failed to delete task : " + err.Error())
		}
	}
}

func TestTask(t *testing.T) {
	suite.Run(t, new(TaskTestSuite))
}

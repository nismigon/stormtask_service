package database

import (
	"naleakan/stormtask/configuration"
	"testing"
)

func BeforeTaskTest() (*DBHandler, int, error) {
	conf, err := configuration.Parse("../configuration.json")
	if err != nil {
		return nil, -1, err
	}
	handler, err := Init(conf.DatabaseURL, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName)
	if err != nil {
		return nil, -1, err
	}
	user, err := handler.AddUser("test@test.com", "Test", "Test", false)
	if err != nil {
		return nil, -1, err
	}
	group, err := handler.AddGroup(user.ID, "TestGroup")
	if err != nil {
		return nil, -1, err
	}
	return handler, group.ID, nil
}

func AfterTaskTest(handler *DBHandler, id int) {
	group, _ := handler.GetGroupByID(id)
	_ = handler.DeleteGroup(id)
	_ = handler.DeleteUser(group.UserID)
	_ = handler.Close()
}

func TestAddAndDeleteTask(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	task, err := handler.AddTask("Task", "A little description", false, false, groupID)
	if err != nil {
		t.Errorf("Failed to add the task : " + err.Error())
	}
	if task.Name != "Task" {
		t.Errorf("Failed to set the name of the task")
	}
	if task.Description != "A little description" {
		t.Errorf("Failed to set the description of the task")
	}
	if task.IsArchived {
		t.Errorf("Failed to set the archived status of the task")
	}
	if task.IsFinished {
		t.Errorf("Failed to set the finished status of the task")
	}
	if task.IDGroup != groupID {
		t.Errorf("Failed to set the owner group of the task")
	}
	err = handler.DeleteTask(task.ID)
	if err != nil {
		t.Errorf("Failed to delete the task : " + err.Error())
	}
	AfterTaskTest(handler, groupID)
}

func TestAddTaskWrongGroup(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	_, err = handler.AddTask("Task", "A little description", false, false, -1)
	if err == nil {
		t.Errorf("Failed to add the task : no error when a wrong group is given")
	}
	AfterTaskTest(handler, groupID)
}

func TestModifyTaskRight(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	tmpTask, err := handler.AddTask("Task", "A little description", false, false, groupID)
	if err != nil {
		t.Errorf("Failed to add the task : " + err.Error())
	}
	task, err := handler.ModifyTask(tmpTask.ID, "New Task", "A bigger description", true, true, tmpTask.IDGroup)
	if err != nil {
		t.Errorf("Failed to modify task : " + err.Error())
	}
	if task.Name != "New Task" {
		t.Errorf("Failed to set the name of the task")
	}
	if task.Description != "A bigger description" {
		t.Errorf("Failed to set the description of the task")
	}
	if !task.IsArchived {
		t.Errorf("Failed to set the archived status of the task")
	}
	if !task.IsFinished {
		t.Errorf("Failed to set the finished status of the task")
	}
	if task.IDGroup != groupID {
		t.Errorf("Failed to set the owner group of the task")
	}
	err = handler.DeleteTask(task.ID)
	if err != nil {
		t.Errorf("Failed to delete the task : " + err.Error())
	}
	AfterTaskTest(handler, groupID)
}

func TestModifyTaskWrongID(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	_, err = handler.ModifyTask(-1, "Test", "Test", false, false, groupID)
	if err == nil {
		t.Errorf("Failed to modify task : give a wrong ID doesn't return an error")
	}
	AfterTaskTest(handler, groupID)
}

func TestModifyTaskWrongGroupID(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	tmpTask, err := handler.AddTask("Test", "Test", false, false, groupID)
	if err != nil {
		t.Errorf("Failed to add task : " + err.Error())
	}
	_, err = handler.ModifyTask(tmpTask.ID, "Test", "Test", false, false, -1)
	if err == nil {
		t.Errorf("Failed to modify task : give a wrong group ID doesn't return an error")
	}
	err = handler.DeleteTask(tmpTask.ID)
	if err != nil {
		t.Errorf("Failed to delete task : " + err.Error())
	}
	AfterTaskTest(handler, groupID)
}

func TestGetTaskByIDRight(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	tmpTask, err := handler.AddTask("Test", "Description", false, false, groupID)
	if err != nil {
		t.Errorf("Failed to add task : " + err.Error())
	}
	task, err := handler.GetTaskByID(tmpTask.ID)
	if err != nil {
		t.Errorf("Failed to get task by its ID : " + err.Error())
	}
	if task.Name != tmpTask.Name {
		t.Errorf("Failed to set the name of the task")
	}
	if task.Description != tmpTask.Description {
		t.Errorf("Failed to set the description of the task")
	}
	if task.IsArchived != tmpTask.IsArchived {
		t.Errorf("Failed to set the archived status of the task")
	}
	if task.IsFinished != tmpTask.IsFinished {
		t.Errorf("Failed to set the finished status of the task")
	}
	if task.IDGroup != tmpTask.IDGroup {
		t.Errorf("Failed to set the owner group of the task")
	}
	err = handler.DeleteTask(tmpTask.ID)
	if err != nil {
		t.Errorf("Failed to delete task : " + err.Error())
	}
	AfterTaskTest(handler, groupID)
}

func TestGetTaskByIDWrongID(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	_, err = handler.GetTaskByID(-1)
	if err == nil {
		t.Errorf("Failed to get task by its ID : No error whereas wrong task ID given")
	}
	AfterTaskTest(handler, groupID)
}

func TestGetTaskByGroupRight(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	tmpTask, err := handler.AddTask("Test", "Description", false, false, groupID)
	if err != nil {
		t.Errorf("Failed to add task : " + err.Error())
	}
	tasks, err := handler.GetTasksByGroup(groupID)
	if err != nil {
		t.Errorf("Failed to get tasks by group : " + err.Error())
	}
	if len(*tasks) != 1 {
		t.Errorf("Failed to get tasks by group : The result of the get task by group is different of 1")
	}
	task := (*tasks)[0]
	if task.Name != tmpTask.Name {
		t.Errorf("Failed to set the name of the task")
	}
	if task.Description != tmpTask.Description {
		t.Errorf("Failed to set the description of the task")
	}
	if task.IsArchived != tmpTask.IsArchived {
		t.Errorf("Failed to set the archived status of the task")
	}
	if task.IsFinished != tmpTask.IsFinished {
		t.Errorf("Failed to set the finished status of the task")
	}
	if task.IDGroup != tmpTask.IDGroup {
		t.Errorf("Failed to set the owner group of the task")
	}
	err = handler.DeleteTask(tmpTask.ID)
	if err != nil {
		t.Errorf("Failed to delete task : " + err.Error())
	}
	AfterTaskTest(handler, groupID)
}

func TestDBHandler_DeleteTasksByGroup(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	_, err = handler.AddTask("Test1", "Description", false, false, groupID)
	if err != nil {
		t.Errorf("Failed to add task : " + err.Error())
	}
	_, err = handler.AddTask("Test2", "Description", false, false, groupID)
	if err != nil {
		t.Errorf("Failed to add task : " + err.Error())
	}
	err = handler.DeleteTasksByGroup(groupID)
	if err != nil {
		t.Errorf("Failed to delete tasks of the group : " + err.Error())
	}
	tasks, err := handler.GetTasksByGroup(groupID)
	if err != nil {
		t.Errorf("Failed to get the tasks for a selected group : " + err.Error())
	}
	if len(*tasks) != 0 {
		t.Errorf("Failed to delete the tasks : The returned array shoul be empty")
	}
	AfterTaskTest(handler, groupID)
}

func TestUniqueConstraintTask(t *testing.T) {
	handler, groupID, err := BeforeTaskTest()
	if err != nil {
		t.Errorf("Failed to initialize the task test, please see other test to see what happens : " + err.Error())
	}
	tmpTask, err := handler.AddTask("Test", "Description", false, false, groupID)
	if err != nil {
		t.Errorf("Failed to add task : " + err.Error())
	}
	task, err := handler.AddTask("Test", "An other description", false, false, groupID)
	if err == nil {
		t.Errorf("Failed to add task : no error whereas break unique constraint on name")
		err = handler.DeleteTask(task.ID)
		if err != nil {
			t.Errorf("Failed to delete task : " + err.Error())
		}
	}
	err = handler.DeleteTask(tmpTask.ID)
	if err != nil {
		t.Errorf("Failed to delete task : " + err.Error())
	}
	AfterTaskTest(handler, groupID)
}

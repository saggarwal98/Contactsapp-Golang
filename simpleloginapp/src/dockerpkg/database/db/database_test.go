package db

import (
	"testing"
	// "github.com/rakyll/gotest"
	"log"
)

var Name, ID []string

func TestAdduser(t *testing.T) {
	status := Adduser("insert into Users(FIRSTNAME,LASTNAME,EMAIL,PASSWORD) VALUES('test','user','testuser@gmail.com','tester');")
	if status == false {
		t.Error("Received error when creating user")
	}
	status = Adduser("insert into Contacts(Name,Phone_number_1,Phone_number_2,Address,email) VALUES('testcontact','1234','5678','testaddress','testuser@gmail.com');")
	if status == false {
		t.Error("Received error when creating contact")
	}
}

func TestCheckUser(t *testing.T) {
	res := Checkuser("testuser@gmail.com", "tester")
	if res != "proceed" && res != "register" && res != "wrong credentials" {
		t.Error("testing failed when checking for user existence in database")
	}
}

func TestGetContacts(t *testing.T) {
	id, name, _, _, _ := GetContacts("testuser@gmail.com")
	if id[0] == "" {
		t.Errorf("testing failed at retrieving contacts for email:testuser@gmail.com")
	} else {
		log.Println(cap(ID))
		ID = make([]string, len(id))
		Name = make([]string, len(name))
		for i, v := range id {
			ID[i] = v
			Name[i] = name[i]
		}
	}
}

func TestGetContactForEdit(t *testing.T) {
	for i, _ := range ID {
		id, name, _, _, _ := GetContactForEdit(ID[i], "testuser@gmail.com")
		if id == "0" {
			if name != "" {
				t.Error("Got error while retrieving contact for edit")
			}
		}
	}
}

func TestGetContactForDelete(t *testing.T) {
	log.Println(ID)
	for i, _ := range ID {
		status := GetContactForDelete(ID[i], Name[i], "testuser@gmail.com")
		if status == 0 {
			t.Errorf("Received error when deleting contact with Name:%s", Name)
		}
	}
}

func TestCleanup(t *testing.T) {
	status := Adduser("Delete from Users where Email='testuser@gmail.com'")
	if status == false {
		t.Error("Received error when deleting testuser")
	}
	status = Adduser("Delete from Contacts where email='testuser@gmail.com'")
	if status == false {
		t.Error("Received error when deleting contacts for testuser@gmail.com")
	}
}

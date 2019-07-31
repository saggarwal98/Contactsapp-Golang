package db
import(
	"testing"
)


func TestAddUser(t *testing.T){
	status:=Adduser("insert into Users(FIRSTNAME,LASTNAME,EMAIL,PASSWORD) VALUES('test','user','testuser@gmail.com','tester');")
	if status == false {
		t.Error("Received error when creating testuser")
	}
	status = Adduser("insert into Contacts(Name,Phone_number_1,Phone_number_2,Address,email) VALUES('testcontact','1234','5678','testaddress','testuser@gmail.com');")
	if status == false {
		t.Error("Received error when creating contact")
	}
	status = Adduser("insert into Contacts(Name,Phone_number_1,Phone_number_2,Address,email) VALUES('','0','0','','testuser@gmail.com');")
	if status == false {
		t.Error("Received error when creating contact")
	}
	status = Adduser("insert into Contacts(ID,Name,Phone_number_1,Phone_number_2,Address,email) VALUES('','','0','0','','');")
	if status == false {
		t.Error("Received error when creating contact with empty id")
	}
	status = Adduser("insert into Contacts(Name,Phone_number_1,Phone_number_2,Address,email) VALUES('','0','0','','')")
	if status == false {
		t.Error("Received error when creating contact with empty email")
	}
	status = Adduser("update Contacts set Address='NewAddress' where Name='testcontact'")
	if status == false {
		t.Error("Received error when updating contact")
	}
	// status = Adduser("delete from Contacts where NOT email='shubhamaggarwalmvn@gmail.com'")
	// if status == false{
	// 	t.Errorf("Received error when deleting contact")
	// }
}



func TestCleanup(t *testing.T){
	status:=Adduser("Delete from Users where EMAIL='testuser@gmail.com'")
	if status == false{
		t.Errorf("Receivded error when deleting user")
	}
}
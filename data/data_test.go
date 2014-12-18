package data

// test data
var users  = []User{
  {
    Name: "Sau Sheong",
    Email: "sausheong@gmail.com",
    Password: "password",
  },
  {
    Name: "John Smith",
    Email: "john@gmail.com",
    Password: "john_pass",
  },
  
}


func setup() {
  ThreadDeleteAll()
  SessionDeleteAll()
  UserDeleteAll()  
}
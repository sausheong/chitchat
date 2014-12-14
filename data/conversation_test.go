package data
//
import (
  "testing"
)

// Delete all conversations from database
func ConversationDeleteAll() (err error) {
  db := db()
  defer db.Close()
  statement := "delete from conversations"
  _, err = db.Exec(statement)
  if err != nil {
    return
  }
  return
}


func Test_CreateConversation(t *testing.T) {
  setup()
  if err := users[0].Create(); err != nil {
    t.Error(err, "Cannot create user.")
  }
  conv, err := users[0].CreateConversation("My first conversation")
  if err != nil {
    t.Error(err, "Cannot create conversation")
  }
  if conv.UserId != users[0].Id {
    t.Error("User not linked with conversation")
  }
}

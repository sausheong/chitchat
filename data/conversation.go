package data

import(
  "time"
)

type Conversation struct {
  Id        int
  Uuid      string
  Topic     string
  UserId    int
  CreatedAt time.Time
}

type Reply struct {
  Id             int
  Uuid           string
  Body           string
  UserId         int
  ConversationId int
  CreatedAt      time.Time
}

// format the CreatedAt date to display nicely on the screen
func (conversation *Conversation) CreatedAtDate() string {
	return conversation.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (reply *Reply) CreatedAtDate() string {
	return reply.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}


// get the number of replies in a conversation
func (conversation *Conversation) NumReplies() (count int) {
  db := db()
  defer db.Close()  
  rows, err := db.Query("SELECT count(*) FROM replies where conversation_id = $1", conversation.Id)
  if err != nil {
    return
  }
  for rows.Next() {
    if err = rows.Scan(&count); err != nil {
      return
    }
  }
  rows.Close()
  return  
}

// get replies to a conversation
func (conversation *Conversation) Replies() (replies []Reply, err error) {
  db := db()
  defer db.Close()  
  rows, err := db.Query("SELECT id, uuid, body, user_id, conversation_id, created_at FROM replies where conversation_id = $1", conversation.Id)
  if err != nil {
    return
  }
  for rows.Next() {
    reply := Reply{}
    if err = rows.Scan(&reply.Id, &reply.Uuid, &reply.Body, &reply.UserId, &reply.ConversationId, &reply.CreatedAt); err != nil {
      return
    }
    replies = append(replies, reply)
  }
  rows.Close()
  return  
}

// Create a new conversation 
func (user *User) CreateConversation(topic string) (conv Conversation, err error) {
  db := db()
  defer db.Close()
  statement := "insert into conversations (uuid, topic, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, topic, user_id, created_at"
  stmt, err := db.Prepare(statement)
  if err != nil {
    return
  }
  defer stmt.Close()
  // use QueryRow to return a row and scan the returned id into the Session struct
  err = stmt.QueryRow(createUUID(), topic, user.Id, time.Now()).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
  if err != nil {
    return
  }    
  return
}



// Create a new reply to a conversation
func (user *User) CreateReply(conv Conversation, body string) (reply Reply, err error) {
  db := db()
  defer db.Close()
  statement := "insert into replies (uuid, body, user_id, conversation_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, conversation_id, created_at"
  stmt, err := db.Prepare(statement)
  if err != nil {
    return
  }
  defer stmt.Close()
  // use QueryRow to return a row and scan the returned id into the Session struct
  err = stmt.QueryRow(createUUID(), body, user.Id, conv.Id, time.Now()).Scan(&reply.Id, &reply.Uuid, &reply.Body, &reply.UserId, &reply.ConversationId, &reply.CreatedAt)
  if err != nil {
    return
  }    
  return
}

// Get all conversations in the database and returns it
func Conversations() (conversations []Conversation, err error){
  db := db()
  defer db.Close()  
  rows, err := db.Query("SELECT id, uuid, topic, user_id, created_at FROM conversations ORDER BY created_at DESC")
  if err != nil {
    return
  }
  for rows.Next() {
    conv := Conversation{}
    if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt); err != nil {
      return
    }
    conversations = append(conversations, conv)
  }
  rows.Close()
  return
}

// Get a conversation by the UUID
func ConversationByUUID(uuid string) (conv Conversation, err error) {
  db := db()
  defer db.Close()  
  conv = Conversation{}
  err = db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM conversations WHERE uuid = $1", uuid).
        Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)  
  return
}


// Get the user who started this conversation
func (conversation *Conversation) User() (user User) {
  db := db()
  defer db.Close()  
  user = User{}
  db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", conversation.UserId).
     Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)  
  return    
}

// Get the user who wrote the reply
func (reply *Reply) User() (user User) {
  db := db()
  defer db.Close()  
  user = User{}
  db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", reply.UserId).
     Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)  
  return    
}

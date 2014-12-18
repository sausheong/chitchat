package data

import(
  "time"
)

type Thread struct {
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
  ThreadId int
  CreatedAt      time.Time
}

// format the CreatedAt date to display nicely on the screen
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (reply *Reply) CreatedAtDate() string {
	return reply.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}


// get the number of replies in a thread
func (thread *Thread) NumReplies() (count int) {
  db := db()
  defer db.Close()  
  rows, err := db.Query("SELECT count(*) FROM replies where thread_id = $1", thread.Id)
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

// get replies to a thread
func (thread *Thread) Replies() (replies []Reply, err error) {
  db := db()
  defer db.Close()  
  rows, err := db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM replies where thread_id = $1", thread.Id)
  if err != nil {
    return
  }
  for rows.Next() {
    reply := Reply{}
    if err = rows.Scan(&reply.Id, &reply.Uuid, &reply.Body, &reply.UserId, &reply.ThreadId, &reply.CreatedAt); err != nil {
      return
    }
    replies = append(replies, reply)
  }
  rows.Close()
  return  
}

// Create a new thread 
func (user *User) CreateThread(topic string) (conv Thread, err error) {
  db := db()
  defer db.Close()
  statement := "insert into threads (uuid, topic, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, topic, user_id, created_at"
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



// Create a new reply to a thread
func (user *User) CreatePost(conv Thread, body string) (reply Reply, err error) {
  db := db()
  defer db.Close()
  statement := "insert into replies (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, thread_id, created_at"
  stmt, err := db.Prepare(statement)
  if err != nil {
    return
  }
  defer stmt.Close()
  // use QueryRow to return a row and scan the returned id into the Session struct
  err = stmt.QueryRow(createUUID(), body, user.Id, conv.Id, time.Now()).Scan(&reply.Id, &reply.Uuid, &reply.Body, &reply.UserId, &reply.ThreadId, &reply.CreatedAt)
  if err != nil {
    return
  }    
  return
}

// Get all threads in the database and returns it
func Threads() (threads []Thread, err error){
  db := db()
  defer db.Close()  
  rows, err := db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
  if err != nil {
    return
  }
  for rows.Next() {
    conv := Thread{}
    if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt); err != nil {
      return
    }
    threads = append(threads, conv)
  }
  rows.Close()
  return
}

// Get a thread by the UUID
func ThreadByUUID(uuid string) (conv Thread, err error) {
  db := db()
  defer db.Close()  
  conv = Thread{}
  err = db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1", uuid).
        Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)  
  return
}


// Get the user who started this thread
func (thread *Thread) User() (user User) {
  db := db()
  defer db.Close()  
  user = User{}
  db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", thread.UserId).
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

package main

import (
  "fmt"
  "net/http"
  "chitchat/data"
)


// GET /conversations/new
// Show the new conversation form page
func newConversation(writer http.ResponseWriter, request *http.Request) {
  loggedin(writer, request)
  t := parseTemplateFiles("layout", "public.navbar", "new.conversation")
  t.Execute(writer, nil)  
}

// POST /signup
// Create the user account
func createConversation(writer http.ResponseWriter, request *http.Request) {
  session := loggedin(writer, request)
  err := request.ParseForm()
  if err != nil {
    danger(err, "Cannot parse form")
  }
  
  user, err := session.User(); if err != nil {
    danger(err, "Cannot get user from session")
  }
  topic := request.PostFormValue("topic")
  if _, err := user.CreateConversation(topic); err != nil {
    danger(err, "Cannot create conversation")
  }
  http.Redirect(writer, request, "/", 302)
}

// GET /conversation/read
// Show the details of the conversation, including the replies and the form to write a reply
func readConversation(writer http.ResponseWriter, request *http.Request) {
  loggedin(writer, request)
  vals := request.URL.Query()
  uuid := vals.Get("id")
  t := parseTemplateFiles("layout", "public.navbar", "conversation")
  conversation, err := data.ConversationByUUID(uuid); if err != nil {
    error_message(writer, request, "Cannot read conversation")
  } else {
    t.Execute(writer, &conversation)  
  }
}

// POST /conversation/reply
// Create the reply
func replyConversation(writer http.ResponseWriter, request *http.Request) {
  session := loggedin(writer, request)
  err := request.ParseForm()
  if err != nil {
    danger(err, "Cannot parse form")
  }
  
  user, err := session.User(); if err != nil {
    danger(err, "Cannot get user from session")
  }
  body := request.PostFormValue("body")
  uuid := request.PostFormValue("uuid")
  conversation, err := data.ConversationByUUID(uuid); if err != nil {
    error_message(writer, request, "Cannot read conversation")
  }
  if _, err := user.CreateReply(conversation, body); err != nil {
    danger(err, "Cannot create reply")
  }
  url := fmt.Sprint("/conversation/read?id=", uuid)
  http.Redirect(writer, request, url , 302)
}
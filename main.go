package main

import (
	"fmt"
	"io"
	"os"
	"log"
	"slices"
	"sync"
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
	"strings"
)

//struct to hold the active sessions
//in case we add future functionalities
var sessions struct{
	sync.RWMutex // thread safety lock
	m map[string]ssh.Session
 }
 //struct to hold the users of rooms
var rooms struct {
	sync.RWMutex
	users map[string][]ssh.Session
}

func broadcast(message string,room string,sender ssh.Session ){
	sessions.RLock()
    defer sessions.RUnlock()

	for _, s := range rooms.users[room]{
		if s == sender{
			io.WriteString(s,fmt.Sprintf("You: %s\n",message))
		}else{
			user_term:=term.NewTerminal(s,"")
			_,err:=user_term.Write([]byte(fmt.Sprintf("\n%s: %s\n%s@localhost:~/%s$ ",sender.User(),message,s.User(),room)))

			if err!=nil{
				fmt.Println("Error writing message: ",err)
				os.Exit(1)
			}
		}
}
}

func main(){

	sessions.m=make(map[string]ssh.Session)
	rooms.users=make(map[string][]ssh.Session)

	rooms.users["A"]=[]ssh.Session{}
	rooms.users["B"]=[]ssh.Session{}
	rooms.users["C"]=[]ssh.Session{}

	ssh.Handle(func (s ssh.Session){
    var cur_room string
    sessionId:=s.Context().SessionID()
	sessions.Lock()
	sessions.m[sessionId]=s
	sessions.Unlock()

	defer func(){
		sessions.Lock()
		delete(sessions.m,sessionId)
		sessions.Unlock()
		fmt.Printf("%s quit the session\n",s.User())
	}()

	user_term:=term.NewTerminal(s,fmt.Sprintf("%s@localhost:~$ ",s.User()))
	
	fmt.Printf("Starting session for %s\n", s.User())

	io.WriteString(s,"Welcome to the chat server\n")
	io.WriteString(s,"Type 'ls' to see the list of commands\n")

	for {
		line, err := user_term.ReadLine()
		if err != nil {
			if err == io.EOF {
				io.WriteString(s, "Session closed\n")
			} else {
				log.Printf("Error reading from session: %v", err)
			}
			break
		}
		line=strings.Trim(line, " ")
		if cur_room=="" {
		if line == "ls" {
			io.WriteString(s, "ls: check the commands\nls -r: check rooms\ncd room_name: enter a room \n:q : exit a room \nexit: exit session\n")
		
			}else if line=="ls -r" && cur_room==""{
			for room,_ :=range rooms.users{
            io.WriteString(s,fmt.Sprintf("%s\n", room))
			}

		}else if len(line)>3 && line[:3]== "cd "{
			room:=line[3:]
			rooms.Lock()
			_,ok:=rooms.users[room]
			rooms.Unlock()
			if !ok{
				io.WriteString(s,"No room found!\n")
			}else{
				rooms.Lock()
				rooms.users[room]=append(rooms.users[room],s)
				rooms.Unlock()
				cur_room=room
				user_term=term.NewTerminal(s,fmt.Sprintf("%s@localhost:~/%s$ ",s.User(),room))
			}
		}else
		if line=="exit" {
			io.WriteString(s,"Goodbye\n")
			break
		}else{
			io.WriteString(s,"Command not found! Try ls to see the commands.\n")
		}
		}else {
			if line==":q"{
				rooms.Lock()
				index:=slices.Index(rooms.users[cur_room],s)
				rooms.users[cur_room]=append(rooms.users[cur_room][:index],rooms.users[cur_room][index+1:]...)
				rooms.Unlock()
				io.WriteString(s,"Weclome to room %s!\n")
				user_term=term.NewTerminal(s,fmt.Sprintf("%s@localhost:~$ ",s.User()))
				cur_room=""
			}else{
			broadcast(line,cur_room,s)
			}
		}
	}
})

   passwordAuth:=ssh.PasswordAuth(func (ctx ssh.Context,pass string)bool{
   return pass=="test123"
   })

	log.Fatal(ssh.ListenAndServe(":3000",nil,
	ssh.HostKeyFile("server_key"),
	passwordAuth))
}
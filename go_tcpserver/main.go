package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"strconv"
)

const (
	connHost = "192.168.1.9"
	connPort = "6666"
	connType = "tcp"
)

type user struct {
	conn   net.Conn
	handle string
	buf    []byte
}

func (u *user) receiveMessage() (string, error) {
	n := 0
	var e error = nil
	for n <= 0 {
		n, e = u.conn.Read(u.buf)
		if e != nil {
			fmt.Println("Error receiving message:", e.Error())
			return "", e
		}
	}

	return string(u.buf[:n]), nil
}

func (u *user) sendMessage(msg string) {
	_, e := u.conn.Write([]byte(msg))
	if e != nil {
		fmt.Println("Error writing:", e.Error())
	}
}

func (u user) killUser() {
	u.conn.Close()
}

type userList struct {
	users []*user
	mtx   *sync.Mutex
}

func (us *userList) findUser(u *user) int {
	for i, e := range us.users {
		if e.conn == u.conn {
			return i
		}
	}

	return -1
}

func (us *userList) getUsers(u *user) string {
	s := strconv.Itoa(len(us.users)) + " people online.\n"
	for _, e := range us.users {
		if e.conn != u.conn {
			s += e.handle + "\n"
		}
	}
	return s
}

func (us *userList) addUser(u *user) bool {
	us.mtx.Lock()
	if len(us.users) < 64 {
		us.users = append(us.users, u)
		us.mtx.Unlock()
		return true
	}
	us.mtx.Unlock()
	u.conn.Write([]byte("Server is currently full... Try again later."))
	return false
}

func (us *userList) removeUser(u *user) {
	i := us.findUser(u)
	if i == -1 {
		fmt.Println("What the fuck...")
		os.Exit(-1)
	}
	(us.users)[i].killUser()

	us.mtx.Lock()
	us.users[len(us.users)-1], us.users[i] = us.users[i], us.users[len(us.users)-1]
	us.users = us.users[:len(us.users)-1]
	us.mtx.Unlock()
}

func (us *userList) broadcast(user *user, msg string) {
	if user != nil {
		msg = user.handle + ": " + msg
	}
	for _, u := range us.users {
		if user == nil || u.conn != user.conn {
			u.sendMessage(msg)
		}
	}
}

func main() {
	fmt.Println("Initializing...")

	lsn, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error while listening:", err.Error())
		os.Exit(-1)
	}

	defer lsn.Close()

	users := new(userList)
	users.mtx = new(sync.Mutex)

	defer shutdown(users)
	
	fmt.Println("Listening on ", connHost, ":", connPort)

	for {

		conn, err := lsn.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(-1)
		}

		go handle(conn, users)

	}
}

func handle(conn net.Conn, users *userList) {
	user := new(user)
	user.buf = make([]byte, 1024)
	user.conn = conn

	_, e := conn.Write([]byte("Enter Handle: \n"))
	if e != nil {
		fmt.Println("Error writing:", e.Error())
		conn.Close()
		return
	}

	handle, e := user.receiveMessage()
	if e != nil {
		conn.Close()
		return
	}

	user.handle = handle[:len(handle)-1]
	users.broadcast(nil, "Welcome "+user.handle+"!\n")

	if !(users.addUser(user)) {
		return
	}
	defer users.removeUser(user)

	for {
		m, e := user.receiveMessage()
		if e == nil {
			switch m {
				case "/users\n":
					user.sendMessage(users.getUsers(user))
				default:
					users.broadcast(user, m)
			}
		} else {
			users.broadcast(nil, user.handle+" has left.\n")
			return
		}
	}
}

func shutdown(users *userList) {
	users.broadcast(nil, "/shutdown")
}

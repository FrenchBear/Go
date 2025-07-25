// g22_test_post05.go
// Learning go, Mastering Go §5 example
// Need to run "go mod tidy" the first time
//
// 2025-06-18	PV		First version

package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/mactsouk/post05"
)

// As the post05 package works with Postgres, there is no need to import lib/pq here:
var MIN = 0
var MAX = 26

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
func getString(length int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == length {
			break
		}
		i++
	}
	return temp
}

// Don't want to commit code with hardcoded password...
func getPassword() string {
	filepath := `C:\Utils\Local\postgres.txt`
	// _, err := os.Stat(filepath)
	// if err != nil {
	// 	panic("Can't find file "+filepath)
	// }
	pwd, err := os.ReadFile(filepath) //
	if err != nil {
		msg := fmt.Sprintf("Error opening %s: %v", filepath, err)
		panic(msg)
	}
	return string(pwd)
}

func main() {
	post05.Hostname = "localhost"
	post05.Port = 5432
	post05.Username = "postgres"
	post05.Password = getPassword()
	post05.Database = "go"
	// This is where you define the connection parameters to the Postgres server, as well
	// as the database you are going to work in (go). As all these variables are in the post05
	// package, they are accessed as such.
	data, err := post05.ListUsers()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range data {
		fmt.Println(v)
	}

	// This is deprecated, no need to call SEED with a random value, it's the default mode. Only call SEED with specific value
	// to get reproducible results.
	// SEED := time.Now().Unix()
	// rand.Seed(SEED)

	random_username := getString(5)

	// We begin by listing existing users.
	// Then, we generate a random string that is used as the username. All randomly
	// generated usernames are 5 characters long because of the getString(5) call. You can
	// change that value if you want.
	t := post05.Userdata{
		Username:    random_username,
		Name:        "Mihalis",
		Surname:     "Tsoukalos",
		Description: "This is me!"}
	id := post05.AddUser(t)
	if id == -1 {
		fmt.Println("There was an error adding user", t.Username)
	}

	err = post05.DeleteUser(id)
	if err != nil {
		fmt.Println(err)
	}
	//Here, we delete the user that we created using the user ID value returned by post05.AddUser(t).
	// Trying to delete it again!
	err = post05.DeleteUser(id)
	if err != nil {
		fmt.Println(err)
	}
	// If you try to delete the same user again, the process fails because the user does not exist.
	id = post05.AddUser(t)
	if id == -1 {
		fmt.Println("There was an error adding user", t.Username)
	}
	// Here, we add the same user again—however, as user ID values are generated by
	// Postgres, this time, the user is going to have a different user ID value than before.
	t = post05.Userdata{
		Username:    random_username,
		Name:        "Mihalis",
		Surname:     "Tsoukalos",
		Description: "This might not be me!"}
	// Here, we update the Description field of the post05.Userdata structure before
	// passing it to post05.UpdateUser(), in order update the information stored in the
	// database.
	err = post05.UpdateUser(t)
	if err != nil {
		fmt.Println(err)
	}
}

package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type User struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (u *User) String() string {
	return fmt.Sprint(string(u.Name))
}

func logAccess(u User) {
	log.Printf("Authenticated: %s", u.String())
}

func AuthenticateCode(code string, users_file string) error {

	// Fetch the list of authorized users.
	users := GetUsers(users_file)

	// Loop through all the users to see if any have a matching code.
	for _, user := range users {
		if code == user.Code {
			logAccess(user)
			return nil
		}
	}

	// If the code doesn't match an entry in the database, return an error.
	return fmt.Errorf("Your code '%s' is invalid, please try again!", code)
}

func GetUsers(users_file string) []User {

	// Check if the file does not exist yet.
	if _, err := os.Stat(users_file); os.IsNotExist(err) {

		// Let user know we created a user file for them.
		log.Print("No users JSON file, creating one now: ", users_file)

		// Create file.
		file, err := os.Create(users_file)
		if err != nil {
			log.Fatal("Error opening/creating users file: ", err.Error())
		}
		defer file.Close()

		// Make the file only readable/writable by the current user.
		err = os.Chmod(users_file, 0600)
		if err != nil {
			log.Fatal("Could not update permissions to user file: ", err.Error())
		}

		// Add an empty JSON list to the file so it can be parsed by the JSON marshaller.
		_, err = file.WriteString("[]")
		if err != nil {
			log.Fatal("Could not add empty JSON hash to file: ", err.Error())
		}

		// Close the file now that we're done with.
		file.Close()
	}

	// Read file contents into a list of bytes
	outputAsBytes, err := ioutil.ReadFile(users_file)
	if err != nil {
		log.Fatal("Error reading file contents: ", err.Error())
	}

	// Parse the JSON in the file and return a slice of User structs.
	var output []User
	err = json.Unmarshal(outputAsBytes, &output)
	if err != nil {
		log.Fatal("Error unmarshalling JSON file: ", err.Error())
	}

	return output
}

// func init() {
// 	log.Print("Hello world")
// }

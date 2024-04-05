package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)


type Person struct {
	name, surname string
	age int
	mail string
	password string
	dateOfAddition time.Time
}


func hashPassword(password string) string {
	hasher := md5.New() // Wiem, że jest słabe do haseł, ale to tylko demonstracyjnie.
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}


func errorCheck(err error) {
	if err != nil {
		fmt.Println("Error during data input:", err)
	}
}

func promptForInput(prompt string, reader *bufio.Reader) string {
    for {
        fmt.Print(prompt + ": ")
        input, err := reader.ReadString('\n')
		errorCheck(err)

		input = strings.TrimSpace(strings.ToLower(input))
        if input == "" {
            fmt.Println("Enter data!")
            continue
        }

        return input
    }
}

func registration (dataBase []Person) []Person {
    var person Person
    reader := bufio.NewReader(os.Stdin)

    name := promptForInput("Enter your Name", reader)
    surname := promptForInput("Enter your Surname", reader)
    ageStr := promptForInput("Enter your Age", reader)
    age, _ := strconv.Atoi(ageStr)
    mail := promptForInput("Enter your E-mail", reader)
    password := promptForInput("Enter your Password", reader)

	// fmt.Printf("Typ danych x: %T\n", ageStr)
	// fmt.Printf("Typ danych x: %T\n", age)
		
	hashedPassword := hashPassword(password)

	person = Person{
		name: name,
		surname: surname,
		age: age, 
		mail: mail,
		password: hashedPassword,
		dateOfAddition: time.Now()}

	dataBase = append(dataBase, person)
	
	return dataBase
}

func handler(w http.ResponseWriter, name, surname string, age int, date time.Time) {
	text := fmt.Sprintf("Hello, %s %s!. You are %d years old. Your account has been created at %s. ", name, surname, age, date.Format("2006-01-02 15:04:05"))
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintf(w, "%s", text)
}

func login(dataBase []Person) {
    reader := bufio.NewReader(os.Stdin)

	mail := promptForInput("Enter your E-mail", reader)
    password := promptForInput("Enter your Password", reader)
    hashedPassword := hashPassword(password)

    for _, person := range dataBase {
        if person.mail == mail {
            if person.password == hashedPassword {
                fmt.Println("Login successful!")

				addr := ":8080"
                http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                    handler(w, person.name, person.surname, person.age, person.dateOfAddition)
                })
				fmt.Printf("Serwer: http://localhost%s\n", addr)

				err := http.ListenAndServe(addr, nil)
				if err != nil {
					fmt.Printf("Błąd uruchamiania serwera: %s\n", err)
				} 

                return
            } else {
                fmt.Println("Incorrect password.")
                return
            }
        }
    }
    fmt.Println("User not found.")
}



func main() {
	var dataBase []Person

	person_1 := Person{name: "Jan", surname: "Kowalski", age: 19, mail: "kowalski@wp.pl", password: hashPassword("kochampieski"), dateOfAddition: time.Now()}
	dataBase = append(dataBase, person_1)

	person_2 := Person{name: "Jan2", surname: "Kowalski2", age: 19, mail: "kowalski2@wp.pl", password: hashPassword("kochampieski2"), dateOfAddition: time.Now()}
	dataBase = append(dataBase, person_2)

	fmt.Println("Sign up: ")
	dataBase = registration(dataBase)

	fmt.Println(dataBase)

	fmt.Println("Sign up: ")
	login(dataBase)

}
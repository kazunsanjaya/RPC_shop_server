package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	Name string
	Price  float64
	Amount float32
}

type VSERVER int

var database []Entry
var dbFileName string

func (a *VSERVER) ListVegetableNames(name string, reply *[]string) error {
	var nameList []string

	for _, val := range database {
		nameList = append(nameList,val.Name)
	}

	*reply = nameList
	return nil
}

func (a *VSERVER) GetPriceByName(name string, reply *float64) error {
	var price float64

	for _, val := range database {
		if val.Name == name {
			price = val.Price
		}
	}
	*reply = price

	return nil
}

func (a *VSERVER) GetAmountByName(name string, reply *float32) error {
	var amount float32

	for _, val := range database {
		if val.Name == name {
			amount = val.Amount
		}
	}
	*reply = amount

	return nil
}

func (a *VSERVER) AddEntry(item Entry, reply *string) error {
	var s bool
	var alreadyAdded bool
	*reply = "Failed to add the Entry!!!"

	alreadyAdded = false
	for _, val := range database {
		if val.Name == item.Name {
			alreadyAdded = true
			break
		}
	}
	if alreadyAdded {
		var res string
		a.UpdateEntry(item, &res)
		if res == "Entry updated successfully!" {
			*reply = "Entry was already available. Entry updated successfully!"
		}
	} else {
		database = append(database, item)
		a.writeDB(dbFileName, &s)
		if s {
			*reply = "Entry added successfully!"
		}
	}

	return nil
}

func (a *VSERVER) loadDB(fileName string, status *bool) error {
	var entryReply string
	var s bool

	f, fErr := os.Open(fileName)
	if fErr != nil {
		log.Fatal("Error in opening vegetableDB.txt")
	}
	defer f.Close()

	dbFileName = fileName
	scanner := bufio.NewScanner(f)
	line := string("")

	fmt.Println("Currently available vegetables in the database: ")
	for scanner.Scan() {
		line = scanner.Text()
		var item = strings.Fields(line)
		price, err3 := strconv.ParseFloat(item[1], 64)
		if err3 != nil {
			log.Fatal("Error occurred when parsing the Price: ", err3)
		}
		amount, err4 := strconv.ParseFloat(item[2], 32)
		if err4 != nil {
			log.Fatal("Error occurred when parsing the Amount: ", err4)
		}
		entry := Entry{item[0], price,float32(amount)}
		a.AddEntry(entry,&entryReply)
		fmt.Println(entry)
	}
	s = true
	*status = s

	return nil
}

func (a *VSERVER) writeDB(fileName string, status *bool) error {

	var s bool
	if fileName != "" {
		f, fErr := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if fErr != nil {
			log.Fatal("Error in opening vegetableDB.txt")
		}

		for _, val := range database {
			f.WriteString(fmt.Sprint(val.Name, " ", val.Price, " ", val.Amount, "\n"))
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	s = true
	*status = s

	return nil
}

func (a *VSERVER) UpdateEntry(item Entry, reply *string) error {
	var s bool
	var found bool
	found =false

	for idx, val := range database {
		if val.Name == item.Name {
			database[idx] = Entry{item.Name, item.Price, item.Amount}
			found = true
			break
		}
	}

	a.writeDB(dbFileName, &s)
	if s && found {
		*reply = "Entry updated successfully!"
	} else {
		*reply = "Failed to update the Entry!!!"
	}
	return nil
}

func main() {

	server := new(VSERVER)
	err := rpc.Register(server)
	if err != nil {
		log.Fatal("Failed to register VSERVER: ", err)
	}

	var dbStatus bool
	server.loadDB("vegetableDB.txt", &dbStatus)
	if dbStatus {
		fmt.Println("DB file read successful...!")
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Error in listening: ", err)
	}
	log.Printf("RPC service started on port %d", 1234)

	http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error in serving: ", err)
	}

}
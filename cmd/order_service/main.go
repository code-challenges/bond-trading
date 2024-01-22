package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

	. "github.com/asalvi0/bond-trading/internal/database"
	. "github.com/asalvi0/bond-trading/internal/models"
	"github.com/asalvi0/bond-trading/internal/order_service"
)

func init() {
	dbURL, err := GetDatabaseURL()
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	db := dbmate.New(u)
	db.SchemaFile = "schema.sql"

	driver, err := db.Driver()
	if err != nil {
		log.Fatal(err)
	}

	exists, err := driver.DatabaseExists()
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		fmt.Println("Initializing Database...")

		err = db.Create()
		if err != nil {
			log.Fatal(err)
		}

		err = db.LoadSchema()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Database Ready!")
	}
}

func main() {
	svc, err := order_service.NewService()
	if err != nil {
		log.Fatal(err)
	}

	svc.ListenForOrders()
}

func newOrders() []*Order {
	return []*Order{
		NewOrder(1, 10, 77, BUY),
		NewOrder(2, 20, 91, BUY),
		NewOrder(3, 30, 90, SELL),
		NewOrder(4, 40, 91, SELL),
		NewOrder(5, 50, 92, SELL),
	}
}

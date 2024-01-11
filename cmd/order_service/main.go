package main

import (
	"embed"
	"fmt"
	"log"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

	"github.com/asalvi0/bond-trading/internal/database"
	. "github.com/asalvi0/bond-trading/internal/database"
	. "github.com/asalvi0/bond-trading/internal/models"
	"github.com/asalvi0/bond-trading/internal/order_service"
)

//go:embed *.sql
var fs embed.FS

func init() {
	dbURL, err := database.URL()
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	db := dbmate.New(u)
	db.FS = fs
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
		db.Create()
		err = db.LoadSchema()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database Ready!")
	}
}

func main() {
	service, err := order_service.NewService()
	if err != nil {
		log.Fatal(err)
	}

	_, err = NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	orders := make([]*Order, 0)

	orders = append(orders, NewOrder(1, 10, 77, BUY))
	orders = append(orders, NewOrder(2, 20, 91, BUY))
	orders = append(orders, NewOrder(3, 30, 90, SELL))
	orders = append(orders, NewOrder(4, 40, 91, SELL))
	orders = append(orders, NewOrder(5, 50, 92, SELL))

	for _, order := range orders {
		err = service.ProcessOrder(order)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(service.PrintOrderBook())
}

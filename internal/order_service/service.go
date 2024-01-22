package order_service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	ob "github.com/i25959341/orderbook"
	"github.com/memphisdev/memphis.go"
	"github.com/shopspring/decimal"

	. "github.com/asalvi0/bond-trading/internal/database"
	"github.com/asalvi0/bond-trading/internal/messaging"
	. "github.com/asalvi0/bond-trading/internal/models"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
)

type Service struct {
	memphisClient *messaging.MemphisClient
	orderBook     *ob.OrderBook
	db            *Database
}

func NewService() (*Service, error) {
	memphisClient, err := messaging.NewMemphisClient()
	if err != nil {
		return nil, err
	}

	orderBook := ob.NewOrderBook()
	if orderBook == nil {
		return nil, errors.New("failed to create order book")
	}

	db, err := NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	svc := Service{
		memphisClient,
		orderBook,
		db,
	}

	memphisClient.SetupConsumers(svc.consumeMessages)

	return &svc, nil
}

func (svc *Service) PrintOrderBook() string { return svc.orderBook.String() }

func (svc *Service) ListenForOrders() {
	// Set up channel on which to send signal notifications
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received
	fmt.Println("Listening for orders...")
	<-c

	// The signal is received, do the cleanup
	svc.db.Close()
	svc.memphisClient.Close()

	fmt.Println("Done!")
}

// Process a batch of messages from a memphis station(topic)
func (svc *Service) consumeMessages(msgs []*memphis.Msg, err error, ctx context.Context) {
	if err != nil {
		log.Printf("Fetch failed: %v", err)
		return
	}

	for i := 0; i < len(msgs); i++ {
		// log.Println(string(msgs[i].Data()))

		order := Order{}
		order.UnmarshalJSON(msgs[i].Data())
		if err != nil {
			log.Println(err)
		}

		err := svc.ProcessOrder(order)
		if err != nil {
			log.Println(err)
		}
		msgs[i].Ack()
	}
}

func (svc *Service) ProcessOrder(order Order) error {
	// write to order book
	done, partial, partialProcessed, err := svc.orderBook.ProcessLimitOrder(
		ob.Side(order.Action.ToSide()),
		order.ID,
		decimal.NewFromInt32(int32(order.Quantity)),
		decimal.NewFromFloat32(order.Price),
	)
	if err != nil {
		return err
	}

	if len(done) > 0 || partial != nil || partialProcessed.GreaterThan(decimal.NewFromInt32(0)) {
		fmt.Printf("Done: %v, Partial: %v, Partial Processed: %v\n", done, partial, partialProcessed)
	}

	// TODO: update orders accordingly
	// if partial != nil {
	// 	order.Quantity = uint(partial.Quantity().IntPart())
	// 	order.Price = float32(partial.Price().InexactFloat64())
	// }

	// if partialProcessed.GreaterThan(decimal.NewFromInt32(0)) {
	// 	order.Quantity = uint(partialProcessed.IntPart())
	// 	order.Price = float32(partialProcessed.InexactFloat64())
	// }

	// svc.db.UpdateOrder(context.Background(), &order)

	return nil
}

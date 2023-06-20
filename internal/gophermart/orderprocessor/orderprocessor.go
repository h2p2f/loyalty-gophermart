package orderprocessor

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"
	"go.uber.org/zap"
	"time"
)

// OrderProcessor is a struct for processing orders
type OrderProcessor struct {
	processor processor
	logger    *zap.Logger
}

// NewOrderProcessor is a constructor for OrderProcessor
func NewOrderProcessor(processor processor, logger *zap.Logger) *OrderProcessor {
	return &OrderProcessor{processor: processor, logger: logger}
}

// Process is a method for processing orders
func (op *OrderProcessor) Process(ctx context.Context, address string) {
	for {
		//Get unfinished orders
		unfinishedOrders, err := op.processor.GetUnfinishedOrders()
		if err != nil {
			return
		}
		//Check status of unfinished orders
		for order, status := range unfinishedOrders {
			op.logger.Sugar().Infof("Order %s status: %s", order, status)
			client := resty.New()
			resp, err := client.R().
				Get(address + "/api/orders/" + order)
			op.logger.Sugar().Infof("Order %s status: %s", order, resp.Status())
			if err != nil {
				return
			}
			if resp.StatusCode() == 204 {
				op.logger.Sugar().Infof("Order %s not found", order)
				continue
			}
			//If we have too many requests, we need to wait
			if resp.StatusCode() == 429 {
				op.logger.Sugar().Infof("Too many requests")
				time.Sleep(60 * time.Second)
			}
			var buf bytes.Buffer
			buf.Write(resp.Body())
			var externalData models.ExternalData
			err = json.Unmarshal(buf.Bytes(), &externalData)
			op.logger.Sugar().Infof("Order %s status: %s", order, externalData.Status)
			if err != nil {
				return
			}
			//Update status of order
			switch externalData.Status {
			//If order is new, we don't need to update status
			case models.NEW:
				if status == models.NEW {
					continue
				}
			//If order is processing...
			case models.PROCESSING:
				//If order was new, we need to update status
				if status == models.NEW {

					err = op.processor.UpdateOrderStatus(order, models.PROCESSING, 0)
					op.logger.Sugar().Infof("Order %s status: %s", order, models.PROCESSING)
					if err != nil {
						return
					}
					//If order was processing, we don't need to update status
				} else {
					continue
				}
			//If order is processed, we need to update status and write accrual to DB
			case models.PROCESSED:
				op.logger.Sugar().Infof("Order %s status: %s", order, models.PROCESSED)

				err = op.processor.UpdateOrderStatus(order, models.PROCESSED, externalData.Accrual)
				if err != nil {
					op.logger.Sugar().Errorf("Error updating order status: %v", err)
					return
				}
			//If order is invalid, we need to update status
			case models.INVALID:
				op.logger.Sugar().Infof("Order %s status: %s", order, models.INVALID)
				err = op.processor.UpdateOrderStatus(order, models.INVALID, 0)
				if err != nil {
					return
				}
			}
		}
		op.logger.Sugar().Infof("Sleeping for 2 seconds")
		//Sleep for 2 seconds
		time.Sleep(2 * time.Second)
	}

}

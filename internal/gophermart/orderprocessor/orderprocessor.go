package orderprocessor

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/logger"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"
	"go.uber.org/zap"
	"time"
)

type OrderProcessor struct {
	processor processor
	logger    *zap.Logger
}

func NewOrderProcessor(processor processor, logger *zap.Logger) *OrderProcessor {
	return &OrderProcessor{processor: processor, logger: logger}
}

func (op *OrderProcessor) Process(address string) {
	for {
		unfinishedOrders, err := op.processor.GetUnfinishedOrders()
		if err != nil {
			return
		}
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
			switch externalData.Status {
			case models.NEW:
				if status == models.NEW {
					continue
				}
			case models.PROCESSING:
				if status == models.NEW {
					err = op.processor.UpdateOrderStatus(order, models.PROCESSING, 0)
					op.logger.Sugar().Infof("Order %s status: %s", order, models.PROCESSING)
					if err != nil {
						return
					}
				} else {
					continue
				}
			case models.PROCESSED:
				op.logger.Sugar().Infof("Order %s status: %s", order, models.PROCESSED)
				err = op.processor.UpdateOrderStatus(order, models.PROCESSED, externalData.Accrual)
				if err != nil {
					logger.Log.Sugar().Errorf("Error updating order %s status: %s", order, models.PROCESSED)
					logger.Log.Sugar().Error(err)
					return
				}
			}
		}
		op.logger.Sugar().Infof("Sleeping for 2 seconds")
		time.Sleep(2 * time.Second)
	}

}

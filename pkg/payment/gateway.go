package payment

import (
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentGateway interface {
	GenerateSnapURL(orderID int, total float64) (string, error)
}

type MidtransClient struct{}

func (m *MidtransClient) GenerateSnapURL(orderID int, totalPrice float64) (string, error) {

	orderIDStr := strconv.Itoa(orderID)

	resp := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderIDStr,
			GrossAmt: int64(totalPrice),
		},
		Expiry: &snap.ExpiryDetails{
			Duration: 15,
			Unit:     "minute",
		},

		Callbacks: &snap.Callbacks{
			Finish: "https://velo-delta-gules.vercel.app",
		},
	}

	snapResp, errMidtrans := snap.CreateTransaction(resp)
	if errMidtrans != nil {
		return "", fmt.Errorf("gagal membuat linkk pembayaran: %v", errMidtrans.GetMessage())
	}

	return snapResp.RedirectURL, nil

}

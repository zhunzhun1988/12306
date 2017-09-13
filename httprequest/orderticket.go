package httprequest

import (
	"12306/log"
	"fmt"
	"net/http"
)

type TicketType int

const (
	Ticket_TDZ = iota
	Ticket_YDZ
	Ticket_EDZ
	Ticket_GJRW
	Ticket_RW
	Ticket_DW
	Ticket_YW
	Ticket_RZ
	Ticket_YZ
	Ticket_WZ
)

func OrderTicket(client *http.Client, secret, date, from, to string) error {
	resp, err := client.PostForm(order_ticket_addr, getOrderTickerUrlValue(secret, date, date, from, to, "dc", "ADULT"))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OrderTicket bad status code:%d", resp.StatusCode)
	}
	body := getBody(resp.Body)
	log.MyOrderLogD("respose body:[%s]", string(body))
	return nil
}

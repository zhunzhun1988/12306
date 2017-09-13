package httprequest

import (
	"12306/log"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func isTicketMatchType(ti *TicketsInfo, types []TicketType) bool {
	if ti.HaveTickets == false {
		return false
	}
	for _, t := range types {
		switch t {
		case Ticket_TDZ:
			if ti.TDZ != "" && ti.TDZ != "无" {
				return true
			}
		case Ticket_YDZ:
			if ti.YDZ != "" && ti.YDZ != "无" {
				return true
			}
		case Ticket_EDZ:
			if ti.EDZ != "" && ti.EDZ != "无" {
				return true
			}
		case Ticket_GJRW:
			if ti.GJRW != "" && ti.GJRW != "无" {
				return true
			}
		case Ticket_RW:
			if ti.RW != "" && ti.RW != "无" {
				return true
			}
		case Ticket_DW:
			if ti.DW != "" && ti.DW != "无" {
				return true
			}
		case Ticket_YW:
			if ti.YW != "" && ti.YW != "无" {
				return true
			}
		case Ticket_RZ:
			if ti.RZ != "" && ti.RZ != "无" {
				return true
			}
		case Ticket_YZ:
			if ti.YZ != "" && ti.YZ != "无" {
				return true
			}
		case Ticket_WZ:
			if ti.WZ != "" && ti.WZ != "无" {
				return true
			}
		}
	}
	return false
}
func OrderTicket(client *http.Client, secret, date, from, to string) error {
	resp, err := client.PostForm(order_ticket_addr, getOrderTickerUrlValue(secret, date, date, from, to, "dc", "ADULT"))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OrderTicket bad status code:%d", resp.StatusCode)
	}
	body := getBody(resp.Body)
	log.MyOrderLogD("OrderTicket:body[%s]", string(body))
	if len(body) == 0 {
		return fmt.Errorf("OrderTicket status is empty")
	}
	otm := OrderTicketMsg{}
	err = json.Unmarshal(body, &otm)
	if err != nil {
		return fmt.Errorf("OrderTicket json Unmarshal err:%v,[%s]", err, string(body))
	}
	if otm.Status == false {
		return fmt.Errorf("OrderTicket fail:%s", strings.Join(otm.Messages, ","))
	}
	return nil
}

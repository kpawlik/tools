package main

import (
	"fmt"
	"time"
)

var (
	i    int
	line = `{"ASSOCIATED_NODE":"NODE_%d","ICOMS_NUMBER":"541043817101","STREET_ADDRESS":"null null null null","DEVICE_TYPE":"power_block","DEVICE_DESCRIPTION":"Power Block (POWER BLOCK)","TICKET_COUNT":"0","IMPACT_NUMBER":"216","VIOLATION_TYPE":"2","VIOLATION_DATE":"2017-11-09 04:53:01.000Z","RULE_DESC":"OUTG","REMEDY":null,"POI_NODE_FAULT_ID":"911201704530193961","POI_ELEMT_FAULT_ID":"604317813yYi13937183321280296","SITE_ID":"541","FAULT_IMPACT_TYPE_ID":"SA","POI_CNT":"2","GNIS_NODE_ID":"941095","XMIT_INDICATOR":"Y","POI_TYPE":"OUTAGE","STATUS":"active"}`
)

func main() {
	ticker := time.NewTicker(time.Millisecond)
	go func(ticker *time.Ticker) {
		for _ = range ticker.C {
			i++
			fmt.Printf("%s\n", fmt.Sprintf(line, i))
		}
	}(ticker)
	time.Sleep(60 * time.Minute)
}

package main

import (
	"context"
	// "fmt"
	"log"
	"time"

	"github.com/ravikantchauhan246/ospy/internal/monitor"
)

func main(){
	checker := monitor.NewChecker(10 * time.Second)

	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://httpstat.us/500",// this will fail
		"https://httpstat.us/200",// This will pass
		"https://ravikant.dev",// This will pass
		
	}

	for _, url := range urls{
		result := checker.Check(context.Background(), url)

		if result.Error != nil {
			log.Printf("❌ %s - Error: %v", url, result.Error)
		}else{
			status := "✅"
			if !result.IsUp {
				status = "❌"
			}
				log.Printf("%s %s - Status: %d, Time: %v", 
				status, url, result.Status, result.ResponseTime)
		}
	}
}


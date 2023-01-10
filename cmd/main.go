package main

/* strava activity webhook subscription callback endpoint manager
   ref: https://developers.strava.com/docs/webhooks/
*/

import (
	"log"

	"github.com/mrpaws/swc/srv"
)

func main() {
	log.Println("strava event subscription service startup")
	srv.Server()
	
}
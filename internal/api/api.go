package api

import (
	"fmt"

	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"gorm.io/gorm"
)

/*
	############## Order of Call Execution ############
						Handler
						   ↓
						Service
						   ↓
						Repository
	###################################################


	############## Dependency Injection ############
						Repository
						   ↓
						 Service
						   ↓
						 Handler
	###################################################
*/

func RunServer(config *config.Config, db *gorm.DB) error {
	router := InitRoutes(db, config)
	port := fmt.Sprintf(":%s", config.ServerPort)
	return router.Run(port)
}

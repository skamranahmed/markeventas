package api

import (
	"fmt"

	"github.com/skamranahmed/markeventas/config"
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

func RunServer(db *gorm.DB) error {
	router := InitRoutes(db)
	port := fmt.Sprintf(":%s", config.ServerPort)
	return router.Run(port)
}

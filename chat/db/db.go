package db

import (
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"../utils"
	"../models"
	log "github.com/sirupsen/logrus"
)

// we need hold a db singleton
// so that db would be connect only once

var (
	// global db connection
	DB *gorm.DB
)

type (
	Todo struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}
)


func ConnectDB(c *utils.AppConfig) bool{
	dbType := c.Database.DbType
	dbUserName := c.Database.DbUserName
	dbUserPassword := c.Database.DbPassword
	dbName := c.Database.DbName

	log.Infof("Database Configuration:    " +
		"db type: %s; " +
			"user: %s; " +
				"db: %s.", dbType, dbUserName, dbName)
	//DB, err := gorm.Open(dbType, dbUserName + ":" + dbUserPassword + "@/" + dbName + "?charset=utf8&parseTime=True&loc=Local&sslmode=disable")
	var err error
	// damn, should not using := here, otherwise DB would be local variable
	DB, err = gorm.Open(dbType, "user=" + dbUserName + " dbname=" + dbName+ " password=" + dbUserPassword + " sslmode=disable")

	utils.CheckError(err, "db.ConnectDB")

	// after connect, I think I should migrating the models if it not exist
	migrateModelsIfNecessary(DB)
	return true
}

// migrate models using gorm AutoMigrate
func migrateModelsIfNecessary(db *gorm.DB){
	// should check tables exist or not, if not then create, otherwise auto-migrate them
	createOrMigrate(db, &models.User{})
	createOrMigrate(db, &models.Group{})
	createOrMigrate(db, &models.Msg{})
}

func createOrMigrate(db *gorm.DB, m interface{}) {
	if db.HasTable(m){
		//log.Infof("%s already exist, auto migrate.", m)
		db.AutoMigrate(m)
	} else {
		// Fuck, somebody tell Create same with CreateTable, shit, not true
		db.CreateTable(m)
	}
}

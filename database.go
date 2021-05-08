package triangulate

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("triangulate.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Println(err)
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&Stat{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&TempSession{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&AuthSession{})
	if err != nil {
		log.Fatal(err)
	}

	stat := Stat{}
	if res := db.First(&stat, "key = ?", "total_generated"); res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Fatal(res.Error)
	}
	if stat.ID == 0 {
		stat.Key = "total_generated"
		stat.Value = 0
		res := db.Create(&stat)
		if res.Error != nil {
			log.Fatal(res.Error)
		}
	}
}

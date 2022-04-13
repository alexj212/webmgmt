package webmgmt

import (
    "fmt"
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    "log"
)

//GetDatabase database connection
func (a *Authr) GetDatabase() *gorm.DB {
    connection, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("database init err: %v", err)
    }
    sqldb, err := connection.DB()
    if err != nil {
        log.Fatalf("database connection err: %v", err)
    }

    err = sqldb.Ping()
    if err != nil {
        log.Fatalf("database ping err: %v", err)
    }
    fmt.Println("Database connection successful.")
    return connection
}

//InitialMigration user table in userdb
func (a *Authr) InitialMigration() {
    connection := a.GetDatabase()
    defer a.CloseDatabase(connection)
    connection.AutoMigrate(User{})
}

//CloseDatabase database connection
func (a *Authr) CloseDatabase(connection *gorm.DB) {
    sqldb, _ := connection.DB()
    sqldb.Close()
}

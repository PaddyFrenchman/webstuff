
package data

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique"`
    Password string
}

func ValidateUser(db *gorm.DB, username, password string) bool {
    var user User
    result := db.Where("username = ? AND password = ?", username, password).First(&user)
    return result.Error == nil
}

func RegisterUser(db *gorm.DB, username, password string) error {
    user := User{Username: username, Password: password}
    return db.Create(&user).Error
}

func GetContacts() []map[string]string {
    return []map[string]string{
        {"name": "Alice", "email": "alice@example.com"},
        {"name": "Bob", "email": "bob@example.com"},
    }
}

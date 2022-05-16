package database

import (
	"fmt"
	"github.com/decadevs/shoparena/models"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// DB provides access to the different db

type DB interface {
	CreateBuyer(user *models.Buyer) (*models.Buyer, error)
	CreateSeller(user *models.Seller) (*models.Seller, error)
	FindAllSellersExcept(except string) ([]models.Seller, error)
	FindBuyerByEmail(email string) (*models.Buyer, error)
	FindBuyerByPhone(phone string) (*models.Buyer, error)
	FindBuyerByUsername(username string) (*models.Buyer, error)
	FindSellerByEmail(email string) (*models.Seller, error)
	FindSellerByPhone(phone string) (*models.Seller, error)
	UpdateUserImageURL(username, url string) error
	FindSellerByUsername(username string) (*models.Seller, error)
	SearchProduct(lowerPrice, upperPrice, category, name string) ([]models.Product, error)
	TokenInBlacklist(token *string) bool
	UpdateBuyerProfile(id uint, update *models.UpdateUser) error
	UpdateSellerProfile(id uint, update *models.UpdateUser) error
	BuyerUpdatePassword(password, newPassword string) (*models.Buyer, error)
	SellerUpdatePassword(password, newPassword string) (*models.Seller, error)
	BuyerResetPassword(email, newPassword string) (*models.Buyer, error)
	CreateBuyerCart(cart *models.Cart) (*models.Cart, error)
	FindIndividualSellerShop(sellerID string) (*models.Seller, error)
	GetAllSellers() ([]models.Seller, error)
	GetProductByID(id string) (*models.Product, error)
}

// Mailer interface to implement mailing service
type Mailer interface {
	SendMail(subject, body, to, Private, Domain string) bool
}

// ValidationError defines error that occur due to validation
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

type DBParams struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
}

func InitDBParams() DBParams {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return DBParams{
		Host:     host,
		User:     user,
		Password: password,
		DbName:   dbName,
		Port:     port,
	}
}

package database

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/decadevs/shoparena/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

//PostgresDb implements the DB interface
type PostgresDb struct {
	DB *gorm.DB
}

// Init sets up the mongodb instance
func (pdb *PostgresDb) Init(host, user, password, dbName, port string) error {
	fmt.Println("connecting to Database.....")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos", host, user, password, dbName, port)
	var err error
	if os.Getenv("DATABASE_URL") != "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if db == nil {
		return fmt.Errorf("database was not initialized")
	} else {
		fmt.Println("Connected to Database")
	}

	err = db.AutoMigrate(&models.Category{}, &models.Seller{}, &models.Product{}, &models.Image{},
		&models.Buyer{}, &models.Cart{}, &models.CartProduct{}, &models.Order{}, &models.Blacklist{})
	if err != nil {
		return fmt.Errorf("migration error: %v", err)
	}

	pdb.DB = db

	return nil

}

// SearchProduct Searches all products from DB
func (pdb *PostgresDb) SearchProduct(lowerPrice, upperPrice, categoryName, name string) ([]models.Product, error) {
	categories := models.Category{}
	var products []models.Product

	LPInt, _ := strconv.Atoi(lowerPrice)
	UPInt, _ := strconv.Atoi(upperPrice)

	if categoryName == "" {
		err := pdb.DB.Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return products, nil
	} else {
		err := pdb.DB.Where("name = ?", categoryName).First(&categories).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	category := categories.ID

	if LPInt == 0 && UPInt == 0 && name == "" {
		err := pdb.DB.Where("category_id = ?", category).Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if LPInt == 0 && name == "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("price <= ?", uint(UPInt)).Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if UPInt == 0 && name == "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("price >= ?", uint(LPInt)).Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if LPInt != 0 && UPInt != 0 && name == "" {
		err := pdb.DB.Where("category_id = ?", category).Where("price >= ?", uint(LPInt)).
			Where("price <= ?", uint(UPInt)).Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if LPInt == 0 && UPInt == 0 && name != "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("title LIKE ?", "%"+name+"%").Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if LPInt == 0 && name != "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("price <= ?", uint(UPInt)).
			Where("title LIKE ?", "%"+name+"%").Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if UPInt == 0 && name != "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("price >= ?", uint(LPInt)).
			Where("title LIKE ?", "%"+name+"%").Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else {
		err := pdb.DB.Where("category_id = ?", category).Where("price >= ?", uint(LPInt)).
			Where("price <= ?", uint(UPInt)).
			Where("title LIKE ?", "%"+name+"%").Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return products, nil
}

// CreateSeller creates a new Seller in the DB
func (pdb *PostgresDb) CreateSeller(user *models.Seller) (*models.Seller, error) {
	var err error
	user.CreatedAt = time.Now()
	user.IsActive = true
	err = pdb.DB.Create(user).Error
	return user, err
}

// CreateBuyer creates a new Buyer in the DB
func (pdb *PostgresDb) CreateBuyer(user *models.Buyer) (*models.Buyer, error) {
	var err error
	user.CreatedAt = time.Now()
	user.IsActive = true
	err = pdb.DB.Create(user).Error
	return user, err
}

//CreateBuyerCart creates a new cart for the buyer
func (pdb *PostgresDb) CreateBuyerCart(cart *models.Cart) (*models.Cart, error) {
	var err error
	cart.CreatedAt = time.Now()
	err = pdb.DB.Create(cart).Error
	return cart, err
}

// FindSellerByUsername finds a user by the username
func (pdb *PostgresDb) FindSellerByUsername(username string) (*models.Seller, error) {
	user := &models.Seller{}

	if err := pdb.DB.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, errors.New("user inactive")
	}
	return user, nil
}

// FindBuyerByUsername finds a user by the username
func (pdb *PostgresDb) FindBuyerByUsername(username string) (*models.Buyer, error) {
	buyer := &models.Buyer{}

	if err := pdb.DB.Where("username = ?", username).First(buyer).Error; err != nil {
		return nil, err
	}
	if !buyer.IsActive {
		return nil, errors.New("user inactive")
	}
	return buyer, nil
}

// FindSellerByEmail finds a user by email
func (pdb *PostgresDb) FindSellerByEmail(email string) (*models.Seller, error) {
	seller := &models.Seller{}
	if err := pdb.DB.Where("email = ?", email).First(seller).Error; err != nil {
		return nil, errors.New(email + " does not exist" + " seller not found")
	}

	return seller, nil
}

// FindBuyerByEmail finds a user by email
func (pdb *PostgresDb) FindBuyerByEmail(email string) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	if err := pdb.DB.Where("email = ?", email).First(buyer).Error; err != nil {
		return nil, errors.New(email + " does not exist" + " buyer not found")
	}

	return buyer, nil
}

// FindSellerByPhone finds a user by the phone
func (pdb PostgresDb) FindSellerByPhone(phone string) (*models.Seller, error) {
	user := &models.Seller{}
	if err := pdb.DB.Where("phone_number =?", phone).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindBuyerByPhone finds a user by the phone
func (pdb PostgresDb) FindBuyerByPhone(phone string) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	if err := pdb.DB.Where("phone_number =?", phone).First(buyer).Error; err != nil {
		return nil, err
	}
	return buyer, nil
}

// TokenInBlacklist checks if token is already in the blacklist collection
func (pdb *PostgresDb) TokenInBlacklist(token *string) bool {
	return false
}

// FindAllUsersExcept returns all the users expcept the one specified in the except parameter
func (pdb *PostgresDb) FindAllSellersExcept(except string) ([]models.Seller, error) {
	sellers := []models.Seller{}
	if err := pdb.DB.Not("username = ?", except).Find(sellers).Error; err != nil {

		return nil, err
	}
	return sellers, nil
}

func (pdb *PostgresDb) UpdateBuyerProfile(id uint, update *models.UpdateUser) error {
	result :=
		pdb.DB.Model(models.Buyer{}).
			Where("id = ?", id).
			Updates(
				models.User{
					FirstName:   update.FirstName,
					LastName:    update.LastName,
					PhoneNumber: update.PhoneNumber,
					Address:     update.Address,
					Email:       update.Email,
				},
			)
	return result.Error
}

func (pdb *PostgresDb) UpdateSellerProfile(id uint, update *models.UpdateUser) error {
	result :=
		pdb.DB.Model(models.Seller{}).
			Where("id = ?", id).
			Updates(
				models.User{
					FirstName:   update.FirstName,
					LastName:    update.LastName,
					PhoneNumber: update.PhoneNumber,
					Address:     update.Address,
					Email:       update.Email,
				},
			)
	return result.Error
}

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func (pdb *PostgresDb) UploadFileToS3(h *session.Session, file multipart.File, fileName string, size int64) (string, error) {
	// get the file size and read the file content into a buffer
	buffer := make([]byte, size)
	file.Read(buffer)
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file you're uploading
	url := "https://s3-eu-west-3.amazonaws.com/arp-rental/" + fileName
	_, err := s3.New(h).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	return url, err
}

func (pdb *PostgresDb) UpdateUserImageURL(username, url string) error {
	result :=
		pdb.DB.Model(models.User{}).
			Where("username = ?", username).
			Updates(
				models.User{
					Image: url,
				},
			)
	return result.Error
}
func (pdb *PostgresDb) BuyerUpdatePassword(password, newPassword string) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	if err := pdb.DB.Model(buyer).Where("password_hash =?", password).Update("password_hash", newPassword).Error; err != nil {
		return nil, err
	}
	return buyer, nil
}
func (pdb *PostgresDb) SellerUpdatePassword(password, newPassword string) (*models.Seller, error) {
	seller := &models.Seller{}
	if err := pdb.DB.Model(seller).Where("password_hash =?", password).Update("password_hash", newPassword).Error; err != nil {
		return nil, err
	}
	return seller, nil
}
func (pdb *PostgresDb) BuyerResetPassword(email, newPassword string) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	if err := pdb.DB.Model(buyer).Where("email =?", email).Update("password_hash", newPassword).Error; err != nil {
		return nil, err
	}
	return buyer, nil
}

//FindIndividualSellerShop return the individual seller and its respective shop gotten by its unique ID
func (pdb *PostgresDb) FindIndividualSellerShop(sellerID string) (*models.Seller, error) {
	//create instance of a seller and its respective product, and unmarshal data into them
	seller := &models.Seller{}

	if err := pdb.DB.Preload("Product").Where("id = ?", sellerID).Find(&seller).Error; err != nil {
		log.Println("Error in finding", err)
		return nil, err
	}

	return seller, nil
}

// GetAllSellers returns all the sellers in the updated database
func (pdb *PostgresDb) GetAllSellers() ([]models.Seller, error) {
	var seller []models.Seller
	err := pdb.DB.Model(&models.Seller{}).Find(&seller).Error
	if err != nil {
		return nil, err
	}
	return seller, nil
}

// GetProductByID returns a particular product by it's ID
func (pdb *PostgresDb) GetProductByID(id string) (*models.Product, error) {
	product := &models.Product{}
	if err := pdb.DB.Where("ID=?", id).First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

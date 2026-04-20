package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"jsonplaceholder-api/internal/config"
	"jsonplaceholder-api/internal/models"
)

// NewConnection creates a new database connection
func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.App.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedUsers seeds the database with JSONPlaceholder users
func SeedUsers(db *gorm.DB) error {
	log.Println("Seeding database with JSONPlaceholder users...")

	// Check if users already exist
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already exist, skipping seeding")
		return nil
	}

	// JSONPlaceholder users data
	users := []models.User{
		{
			ID:       1,
			Name:     "Leanne Graham",
			Username: "Bret",
			Email:    "Sincere@april.biz",
			Phone:    "1-770-736-8031 x56442",
			Website:  "hildegard.org",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Kulas Light",
				Suite:   "Apt. 556",
				City:    "Gwenborough",
				Zipcode: "92998-3874",
				Geo: models.Geo{
					Lat: "-37.3159",
					Lng: "81.1496",
				},
			},
			Company: models.Company{
				Name:        "Romaguera-Crona",
				CatchPhrase: "Multi-layered client-server neural-net",
				BS:          "harness real-time e-markets",
			},
		},
		{
			ID:       2,
			Name:     "Ervin Howell",
			Username: "Antonette",
			Email:    "Shanna@melissa.tv",
			Phone:    "010-692-6593 x09125",
			Website:  "anastasia.net",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Victor Plains",
				Suite:   "Suite 879",
				City:    "Wisokyburgh",
				Zipcode: "90566-7771",
				Geo: models.Geo{
					Lat: "-43.9509",
					Lng: "-34.4618",
				},
			},
			Company: models.Company{
				Name:        "Deckow-Crist",
				CatchPhrase: "Proactive didactic contingency",
				BS:          "synergize scalable supply-chains",
			},
		},
		{
			ID:       3,
			Name:     "Clementine Bauch",
			Username: "Samantha",
			Email:    "Nathan@yesenia.net",
			Phone:    "1-463-123-4447",
			Website:  "ramiro.info",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Douglas Extension",
				Suite:   "Suite 847",
				City:    "McKenziehaven",
				Zipcode: "59590-4157",
				Geo: models.Geo{
					Lat: "-68.6102",
					Lng: "-47.0653",
				},
			},
			Company: models.Company{
				Name:        "Romaguera-Jacobson",
				CatchPhrase: "Face to face bifurcated interface",
				BS:          "e-enable strategic applications",
			},
		},
		{
			ID:       4,
			Name:     "Patricia Lebsack",
			Username: "Karianne",
			Email:    "Julianne.OConner@kory.org",
			Phone:    "493-170-9623 x156",
			Website:  "kale.biz",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Hoeger Mall",
				Suite:   "Apt. 692",
				City:    "South Elvis",
				Zipcode: "53919-4257",
				Geo: models.Geo{
					Lat: "29.4572",
					Lng: "-164.2990",
				},
			},
			Company: models.Company{
				Name:        "Robel-Corkery",
				CatchPhrase: "Multi-tiered zero tolerance productivity",
				BS:          "transition cutting-edge web services",
			},
		},
		{
			ID:       5,
			Name:     "Chelsey Dietrich",
			Username: "Kamren",
			Email:    "Lucio_Hettinger@annie.ca",
			Phone:    "(254)954-1289",
			Website:  "demarco.info",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Skiles Walks",
				Suite:   "Suite 351",
				City:    "Roscoeview",
				Zipcode: "33263",
				Geo: models.Geo{
					Lat: "-31.8129",
					Lng: "62.5342",
				},
			},
			Company: models.Company{
				Name:        "Keebler LLC",
				CatchPhrase: "User-centric fault-tolerant solution",
				BS:          "revolutionize end-to-end systems",
			},
		},
		{
			ID:       6,
			Name:     "Mrs. Dennis Schulist",
			Username: "Leopoldo_Corkery",
			Email:    "Karley_Dach@jasper.info",
			Phone:    "1-477-935-8478 x6430",
			Website:  "ola.org",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Norberto Crossing",
				Suite:   "Apt. 950",
				City:    "South Christy",
				Zipcode: "23505-1337",
				Geo: models.Geo{
					Lat: "-71.4197",
					Lng: "71.7478",
				},
			},
			Company: models.Company{
				Name:        "Considine-Lockman",
				CatchPhrase: "Synchronised bottom-line interface",
				BS:          "e-enable innovative applications",
			},
		},
		{
			ID:       7,
			Name:     "Kurtis Weissnat",
			Username: "Elwyn.Skiles",
			Email:    "Telly.Hoeger@billy.biz",
			Phone:    "210.067.6132",
			Website:  "elvis.io",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Rex Trail",
				Suite:   "Suite 280",
				City:    "Howemouth",
				Zipcode: "58804-1099",
				Geo: models.Geo{
					Lat: "24.8918",
					Lng: "21.8984",
				},
			},
			Company: models.Company{
				Name:        "Johns Group",
				CatchPhrase: "Configurable multimedia task-force",
				BS:          "generate enterprise e-tailers",
			},
		},
		{
			ID:       8,
			Name:     "Nicholas Runolfsdottir V",
			Username: "Maxime_Nienow",
			Email:    "Sherwood@rosamond.me",
			Phone:    "586.493.6943 x140",
			Website:  "jacynthe.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Ellsworth Summit",
				Suite:   "Suite 729",
				City:    "Aliyaview",
				Zipcode: "45169",
				Geo: models.Geo{
					Lat: "-14.3990",
					Lng: "-120.7677",
				},
			},
			Company: models.Company{
				Name:        "Abernathy Group",
				CatchPhrase: "Implemented secondary concept",
				BS:          "e-enable extensible e-tailers",
			},
		},
		{
			ID:       9,
			Name:     "Glenna Reichert",
			Username: "Delphine",
			Email:    "Chaim_McDermott@dana.io",
			Phone:    "(775)976-6794 x41206",
			Website:  "conrad.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Dayna Park",
				Suite:   "Suite 449",
				City:    "Bartholomebury",
				Zipcode: "76495-3109",
				Geo: models.Geo{
					Lat: "24.6463",
					Lng: "-168.8889",
				},
			},
			Company: models.Company{
				Name:        "Yost and Sons",
				CatchPhrase: "Switchable contextually-based project",
				BS:          "aggregate real-time technologies",
			},
		},
		{
			ID:       10,
			Name:     "Clementina DuBuque",
			Username: "Moriah.Stanton",
			Email:    "Rey.Padberg@karina.biz",
			Phone:    "024-648-3804",
			Website:  "ambrose.net",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Address: models.Address{
				Street:  "Kattie Turnpike",
				Suite:   "Suite 198",
				City:    "Lebsackbury",
				Zipcode: "31428-2261",
				Geo: models.Geo{
					Lat: "-38.2386",
					Lng: "57.2232",
				},
			},
			Company: models.Company{
				Name:        "Hoeger LLC",
				CatchPhrase: "Centralized empowering task-force",
				BS:          "target end-to-end models",
			},
		},
	}

	// Create users in batches
	for _, user := range users {
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Username, err)
			continue
		}
	}

	log.Printf("Successfully seeded %d users", len(users))
	return nil
} 
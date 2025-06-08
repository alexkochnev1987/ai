package database

import (
	"context"
	"fmt"
	"time"

	"jsonplaceholder-api/internal/config"
	"jsonplaceholder-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database wraps the GORM database connection
type Database struct {
	*gorm.DB
}

// New creates a new database connection
func New(cfg *config.Config) (*Database, error) {
	// Configure GORM logger based on environment
	var gormLogger logger.Interface
	if cfg.IsDevelopment() {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN()), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
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
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{DB: db}, nil
}

// AutoMigrate runs database migrations
func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	)
}

// Seed populates the database with initial data
func (d *Database) Seed() error {
	// Check if users already exist
	var count int64
	if err := d.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 {
		return nil // Already seeded
	}

	// Seed users with data from JSONPlaceholder
	users := []models.User{
		{
			Name:     "Leanne Graham",
			Username: "Bret",
			Email:    "Sincere@april.biz",
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
			Phone:   "1-770-736-8031 x56442",
			Website: "hildegard.org",
			Company: models.Company{
				Name:        "Romaguera-Crona",
				CatchPhrase: "Multi-layered client-server neural-net",
				BS:          "harness real-time e-markets",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Ervin Howell",
			Username: "Antonette",
			Email:    "Shanna@melissa.tv",
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
			Phone:   "010-692-6593 x09125",
			Website: "anastasia.net",
			Company: models.Company{
				Name:        "Deckow-Crist",
				CatchPhrase: "Proactive didactic contingency",
				BS:          "synergize scalable supply-chains",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Clementine Bauch",
			Username: "Samantha",
			Email:    "Nathan@yesenia.net",
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
			Phone:   "1-463-123-4447",
			Website: "ramiro.info",
			Company: models.Company{
				Name:        "Romaguera-Jacobson",
				CatchPhrase: "Face to face bifurcated interface",
				BS:          "e-enable strategic applications",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Patricia Lebsack",
			Username: "Karianne",
			Email:    "Julianne.OConner@kory.org",
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
			Phone:   "493-170-9623 x156",
			Website: "kale.biz",
			Company: models.Company{
				Name:        "Robel-Corkery",
				CatchPhrase: "Multi-tiered zero tolerance productivity",
				BS:          "transition cutting-edge web services",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Chelsey Dietrich",
			Username: "Kamren",
			Email:    "Lucio_Hettinger@annie.ca",
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
			Phone:   "(254)954-1289",
			Website: "demarco.info",
			Company: models.Company{
				Name:        "Keebler LLC",
				CatchPhrase: "User-centric fault-tolerant solution",
				BS:          "revolutionize end-to-end systems",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Mrs. Dennis Schulist",
			Username: "Leopoldo_Corkery",
			Email:    "Karley_Dach@jasper.info",
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
			Phone:   "1-477-935-8478 x6430",
			Website: "ola.org",
			Company: models.Company{
				Name:        "Considine-Lockman",
				CatchPhrase: "Synchronised bottom-line interface",
				BS:          "e-enable innovative applications",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Kurtis Weissnat",
			Username: "Elwyn.Skiles",
			Email:    "Telly.Hoeger@billy.biz",
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
			Phone:   "210.067.6132",
			Website: "elvis.io",
			Company: models.Company{
				Name:        "Johns Group",
				CatchPhrase: "Configurable multimedia task-force",
				BS:          "generate enterprise e-tailers",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Nicholas Runolfsdottir V",
			Username: "Maxime_Nienow",
			Email:    "Sherwood@rosamond.me",
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
			Phone:   "586.493.6943 x140",
			Website: "jacynthe.com",
			Company: models.Company{
				Name:        "Abernathy Group",
				CatchPhrase: "Implemented secondary concept",
				BS:          "e-enable extensible e-tailers",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Glenna Reichert",
			Username: "Delphine",
			Email:    "Chaim_McDermott@dana.io",
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
			Phone:   "(775)976-6794 x41206",
			Website: "conrad.com",
			Company: models.Company{
				Name:        "Yost and Sons",
				CatchPhrase: "Switchable contextually-based project",
				BS:          "aggregate real-time technologies",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
		{
			Name:     "Clementina DuBuque",
			Username: "Moriah.Stanton",
			Email:    "Rey.Padberg@karina.biz",
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
			Phone:   "024-648-3804",
			Website: "ambrose.net",
			Company: models.Company{
				Name:        "Hoeger LLC",
				CatchPhrase: "Centralized empowering task-force",
				BS:          "target end-to-end models",
			},
			PasswordHash: "$2a$10$X8Y1Z2A3B4C5D6E7F8G9H0", // bcrypt hash of "password123"
		},
	}

	// Create users in batch
	if err := d.DB.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health checks database health
func (d *Database) Health(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
} 
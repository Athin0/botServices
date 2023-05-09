package db

import (
	"botServices/pkg/model"
	"database/sql"
	_ "encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type PostgresDB struct {
	Client *sql.DB
}

func NewPostgresDB(cfg Config) (*PostgresDB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("postgres connect error : (%v)", err)
	}
	fmt.Println(db)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{Client: db}, nil
}

func InitDB() (*PostgresDB, error) {
	viper.AddConfigPath("../botService/db") //change

	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("error in reading config: %v", err)
		return nil, err
	}
	db, err := NewPostgresDB(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("error creating db: %v \n", err)
		return nil, err
	}
	return db, nil
}

func (db *PostgresDB) Set(owner int64, login string, password string, service string) (int, error) {
	queryRow := db.Client.QueryRow(
		"INSERT INTO Services(service,login, password,user_id) values ($1, $2, $3,$4)",
		service, login, password, owner,
	)
	if queryRow.Err() != nil {
		log.Printf(queryRow.Err().Error())
		return 0, queryRow.Err()
	}
	return 0, nil
}

func (db *PostgresDB) Get(owner int64, service string) (*model.ServiceInfo, error) {
	err := db.Client.QueryRow(
		"SELECT service, login, password FROM Services WHERE user_id = $1 and service =$2",
		owner, service,
	)
	if err.Err() != nil {
		log.Printf(err.Err().Error())
		return nil, err.Err()
	}
	var ans model.ServiceInfo
	err2 := err.Scan(&ans.Service, &ans.Login, &ans.Password)
	if err2 != nil {
		return nil, err2
	}
	return &ans, nil
}

func (db *PostgresDB) Del(owner int64, service string) error {
	err := db.Client.QueryRow(
		"DELETE FROM Services WHERE user_id = $1 and service =$2",
		owner, service,
	)
	if err.Err() != nil {
		log.Printf(err.Err().Error())
		return err.Err()
	}
	return nil
}

func (db *PostgresDB) GetAll(owner int64) ([]model.ServiceInfo, error) {
	rows, err := db.Client.Query(
		"SELECT service, login, password FROM Services WHERE user_id = $1 ",
		owner,
	)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	var ans []model.ServiceInfo
	var a model.ServiceInfo
	for rows.Next() {
		err2 := rows.Scan(&a.Service, &a.Login, &a.Password)
		if err2 != nil {
			return nil, err2
		}
		ans = append(ans, a)
	}
	return ans, nil
}

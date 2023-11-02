package postgresql

import (
	"database/sql"
	"log"
	"rec/internal/config"
	"rec/internal/pkg/postgresql/model"

	_ "github.com/lib/pq"
)

type DBService struct {
	db *sql.DB
}

func NewDB(database *sql.DB) *DBService {
	return &DBService{db: database}
}

func Conn(cfg config.Config) (*DBService, error) {
	var err error
	dbConn := DBService{}

	dsn := "user=" + cfg.Storage.Username + " password=" + cfg.Storage.Password + " dbname=" + cfg.Storage.Database + " sslmode=disable"
	dbConn.db, err = sql.Open(cfg.Storage.DriverName, dsn)
	if err != nil {
		return &DBService{}, err
	}
	return &dbConn, err
}

func (s *DBService) Close() error {
	err := s.db.Close()
	return err
}

func (s *DBService) CreateOrder(jsonData *model.OrderItem) (sql.Result, error) {
	q, err := s.db.Exec(`insert into orders(id, orderdata) values ($1, $2)`, jsonData.ID, jsonData.Data)
	if err != nil {
		log.Println("Ошибка при создании нового заказа, возможно он уже существует")
	}
	return q, err
}

func (s *DBService) Orders() ([]model.OrderItem, error) {
	rows, err := s.db.Query("select * from orders")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rowItem := model.OrderItem{}
	rows.Scan(&rowItem.ID, &rowItem.Data)
	defer rows.Close()
	orders := []model.OrderItem{}
	for rows.Next() {
		str := model.OrderItem{}
		err := rows.Scan(&str.ID, &str.Data)
		if err != nil {
			return orders, err
		}
		orders = append(orders, str)
	}
	return orders, err
}

func (s *DBService) OrderById(id string) (*model.OrderItem, error) {
	row := s.db.QueryRow("select * from orders where id=$1", id)
	rowData := new(model.OrderItem)
	err := row.Scan(&rowData.ID, &rowData.Data)
	return rowData, err
}

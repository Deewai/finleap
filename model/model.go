package model

import (
	"database/sql"
	"fmt"
)

type Conn struct {
	DbHost string
	DbPort string
	DbName string
	DbUser string
	DbPass string
}

type City struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Temperature struct {
	ID        int   `json:"id"`
	CityID    int   `json:"city_id"`
	Max       int   `json:"max"`
	Min       int   `json:"min"`
	Timestamp int64 `json:"timestamp"`
}

type Forecast struct {
	CityID int     `json:"city_id"`
	Max    float32 `json:"max"`
	Min    float32 `json:"min"`
	Sample int     `json:"sample"`
}

type Webhook struct {
	ID          int    `json:"id"`
	CityID      int    `json:"city_id"`
	CallbackURL string `json:"callback_url"`
}

func NewConn(protocol, host, port, user, password, dbname string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", user, password, host, port, dbname)

	db, err := sql.Open(protocol, connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *City) Create(db *sql.DB) error {
	sql := fmt.Sprintf("INSERT INTO cities(name, latitude, longitude) VALUES('%s', %f, %f)", c.Name, c.Latitude, c.Longitude)
	res, err := db.Exec(sql)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return nil
}

func (c *City) get(db *sql.DB) error {
	sql := fmt.Sprintf("SELECT name, latitude, longitude FROM cities WHERE id=%d", c.ID)
	return db.QueryRow(sql).Scan(&c.Name, &c.Latitude, &c.Longitude)
}

func (c *City) Update(db *sql.DB) error {
	sql := fmt.Sprintf("UPDATE cities SET name='%s', latitude=%f, longitude=%f WHERE id=%d", c.Name, c.Latitude, c.Longitude, c.ID)
	_, err := db.Exec(sql)
	return err
}

func (c *City) Delete(db *sql.DB) error {
	err := c.get(db)
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM cities WHERE id=%d", c.ID)
	_, err = db.Exec(sql)
	return err
}

func (t *Temperature) Create(db *sql.DB) error {
	sql := fmt.Sprintf("INSERT INTO temperatures(city_id, max, min, timestamp) VALUES('%d', %d, %d, FROM_UNIXTIME(%d))", t.CityID, t.Max, t.Min, t.Timestamp)
	res, err := db.Exec(sql)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = int(id)
	return nil
}

func GetTemperatures(db *sql.DB, CityID int, timestamp int64) ([]Temperature, error) {
	sql := fmt.Sprintf("SELECT id, city_id, max, min FROM temperatures WHERE city_id = %d AND timestamp >= FROM_UNIXTIME(%d)", CityID, timestamp)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	temperatures := []Temperature{}
	for rows.Next() {
		var t Temperature
		if err := rows.Scan(&t.ID, &t.CityID, &t.Max, &t.Min); err != nil {
			return nil, err
		}
		temperatures = append(temperatures, t)
	}
	return temperatures, nil
}

func (w *Webhook) Create(db *sql.DB) error {
	sql := fmt.Sprintf("INSERT INTO webhooks(city_id, callback_url) VALUES(%d, '%s')", w.CityID, w.CallbackURL)
	res, err := db.Exec(sql)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	w.ID = int(id)
	return nil
}

func GetWebhooks(db *sql.DB) ([]Webhook, error) {
	sql := "SELECT * FROM webhooks"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	webhooks := []Webhook{}
	for rows.Next() {
		var w Webhook
		if err := rows.Scan(&w.ID, &w.CityID, &w.CallbackURL); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, nil
}

func (w *Webhook) get(db *sql.DB) error {
	sql := fmt.Sprintf("SELECT id, city_id, callback_url FROM webhooks WHERE id=%d", w.ID)
	return db.QueryRow(sql).Scan(&w.ID, &w.CityID, &w.CallbackURL)
}

func (w *Webhook) Delete(db *sql.DB) error {
	err := w.get(db)
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM webhooks WHERE id=%d", w.ID)
	_, err = db.Exec(sql)
	return err
}

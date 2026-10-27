package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) "+
		"VALUES (:client, :status, :address, :created_at);",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return -1, err
	}

	answer, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(answer), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	res := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number; ",
		sql.Named("number", number))

	p := Parcel{}
	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	row, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client;",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}

	var res []Parcel
	for row.Next() {
		p := Parcel{}
		err = row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err = row.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number;",
		sql.Named("number", number),
		sql.Named("status", status))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	status := ""
	res := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number;",
		sql.Named("number", number))
	err := res.Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("ERROR: status is not registered")
	}

	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number;",
		sql.Named("address", address),
		sql.Named("number", number))

	return nil
}

func (s ParcelStore) Delete(number int) error {
	status := ""
	res := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number;",
		sql.Named("number", number))
	err := res.Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("ERROR: status is not registered")
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number;",
		sql.Named("number", number))

	return nil
}

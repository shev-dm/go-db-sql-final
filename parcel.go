package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client,status,address,created_at) "+
		"values (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	index, _ := res.LastInsertId()

	return int(index), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	row := s.db.QueryRow("SELECT * from parcel where number = :number",
		sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel
	rows, err := s.db.Query("SELECT * from parcel where client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var parcel Parcel
		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, parcel)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel set status = :status where number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number)

	if err != nil {
		return err
	}

	if parcel.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("UPDATE parcel set address = :address where number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	parcel, err := s.Get(number)

	if err != nil {
		return err
	}

	if parcel.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

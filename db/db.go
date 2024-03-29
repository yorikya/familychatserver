package db

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

//DataBase app data database
type DataBase struct {
	db              *bolt.DB
	filePath        string
	usersAuthBucket string
}

func (d *DataBase) createBucketTX(tx *bolt.Tx, name string) (*bolt.Bucket, error) {
	b, err := tx.CreateBucket([]byte(name))
	if err != nil {
		return nil, fmt.Errorf("db: failed to create '%s' bucket", name)
	}
	return b, nil
}

func (d *DataBase) createAuthUsers() error {
	err := d.CreateAuthUser("yuri", "a")
	if err != nil {
		return err
	}

	err = d.CreateAuthUser("dani", "123456")
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) CreateBucket(name string) error {
	//Create bucket
	return d.db.Update(func(tx *bolt.Tx) error {
		// authentication bucket
		_, err := d.createBucketTX(tx, name)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *DataBase) init() {
	//Create buckets
	d.db.Update(func(tx *bolt.Tx) error {
		// authentication bucket
		_, err := d.createBucketTX(tx, d.usersAuthBucket)
		if err != nil {
			return err
		}

		return nil
	})

	//Create users
	err := d.createAuthUsers()
	if err != nil {
		panic(err)
	}

}

//CreateAuthUser create a new user in auth table
func (d *DataBase) CreateAuthUser(user, pass string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.usersAuthBucket))
		err := b.Put([]byte(user), []byte(pass))
		return err
	})
}

//GetAuthUser re turn table vlaue of auth user
func (d *DataBase) GetAuthUser(name string) (val string, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.usersAuthBucket))
		bVal := b.Get([]byte(name))
		if bVal == nil {
			return fmt.Errorf("the user: '%s' does not exists", name)
		}

		val = string(bVal)
		return nil
	})
	return
}

func (d *DataBase) Close() {
	d.db.Close()
}

//NewDB starts a new DB
func NewDB() *DataBase {
	filePath := "db/bolt.db"
	boltDB, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	database := &DataBase{
		db:              boltDB,
		filePath:        filePath,
		usersAuthBucket: "usersAuth",
	}

	database.init()
	return database

}

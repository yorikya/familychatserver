package db

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

//DataBase app data database
type DataBase struct {
	db                *bolt.DB
	filePath          string
	usersAuthBucket   string
	usersRoomsBuckets string
}

func (d *DataBase) createBucket(tx *bolt.Tx, name string) (*bolt.Bucket, error) {
	b, err := tx.CreateBucket([]byte(name))
	if err != nil {
		return nil, fmt.Errorf("db: failed to create '%s' bucket", name)
	}
	return b, nil
}

func (d *DataBase) makeRoomBucket(tx *bolt.Tx, roomID string) error {
	_, err := d.createBucket(tx, fmt.Sprintf("%s_%s", d.usersRoomsBuckets, roomID))
	if err != nil {
		return err
	}
	return nil
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

func (d *DataBase) init() {
	//Create buckets
	d.db.Update(func(tx *bolt.Tx) error {
		// authentication bucket
		_, err := d.createBucket(tx, d.usersAuthBucket)
		if err != nil {
			return err
		}

		// rooms buckets
		err = d.makeRoomBucket(tx, "1")
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

//NewDB starts a new DB
func NewDB() *DataBase {
	filePath := "db/bolt.db"
	boltDB, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	database := &DataBase{
		db:                boltDB,
		filePath:          filePath,
		usersRoomsBuckets: "usersRoom",
		usersAuthBucket:   "usersAuth",
	}

	database.init()
	return database

}

package db

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

//DataBase app data database
type DataBase struct {
	db                  *bolt.DB
	filePath            string
	usersAuthBucket     string
	usersRoomsBuckets   string
	usersRoomLogBuckets string
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

func (d *DataBase) makeRoomLogBacket(tx *bolt.Tx, roomID string) error {
	_, err := d.createBucket(tx, fmt.Sprintf("%s_%s", d.usersRoomLogBuckets, roomID))
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) makeRoom(tx *bolt.Tx, roomID string) error {
	err := d.makeRoomBucket(tx, roomID)
	if err != nil {
		return err
	}

	err = d.makeRoomLogBacket(tx, roomID)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) createAuthUsers() error {
	err := d.CreateAuthUser("yuri", "123456")
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
		err = d.makeRoom(tx, "1")
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

//AddChatLog
func (d *DataBase) AddChatLog(msgID, msg string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf("%s_1", d.usersRoomLogBuckets)))
		err := b.Put([]byte(msgID), []byte(msg))
		return err
	})
}

//GetChatLog
func (d *DataBase) GetChatLog(msgID string) (val string, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf("%s_1", d.usersRoomLogBuckets)))
		bVal := b.Get([]byte(msgID))
		if bVal == nil {
			return fmt.Errorf("the messageID: '%s' does not exists", msgID)
		}

		val = string(bVal)
		return nil
	})
	return
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
		db:                  boltDB,
		filePath:            filePath,
		usersRoomsBuckets:   "usersRoom",
		usersRoomLogBuckets: "usersLog",
		usersAuthBucket:     "usersAuth",
	}

	database.init()
	return database

}

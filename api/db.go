package main

import (
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

////////////////
// DATA TYPES //
////////////////

type User struct {
	ID       int64
	Username string
	Password string
	Admin    bool
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %s %t>", u.ID, u.Username, u.Password, u.Admin)
}

type Address struct {
	ID      int64
	Address string
	UserID  *User
}

func (a Address) String() string {
	return fmt.Sprintf("Address<%d %s>", a.ID, a.Address)
}

type Tree struct {
	ID     int64
	Status bool
	Health int8
}

func (t Tree) String() string {
	return fmt.Sprintf("Tree<%d %t %d>", t.ID, t.Status, t.Health)
}

///////////////////////
// DATABASE MANAGERS //
///////////////////////

func getConnectOpts() *pg.Options {
	return &pg.Options{
		// With default settings Postgres creates a database with a host of localhost and a role
		// inherited from the system user used to create it. I am "alekseybilogur" on my work
		// machine, so I am "alekseybilogur" in the database as well.
		//
		// The default database created is named "Postgres". It is possible to create a new
		// database (CREATE DATABASE), switch to it (\connect), and drop the Postgres one.
		// Alternatively you may initialize a database under a different name at the outset.
		//
		// To list databases on a certain host from the command line:
		// psql -h localhost --username=alekseybilogur --list
		User:            "alekseybilogur",
		Addr:            "localhost:5432",
		Database:        "trees",
		ApplicationName: "trees",
	}
}

func createSchema(db *pg.DB) error {
	//  This next line takes some work to understand. We are ranging over a slice consisting of
	// anonymous interfaces (e.g. any-type structs).
	//
	// Each entry in the slice is a value reference of one of types defined by the ORM mapping.
	// Creating the value reference requires specifying fields. In this case nil is a special
	// sentinal value telling the program to use all-default values.
	for _, model := range []interface{}{(*User)(nil), (*Address)(nil), (*Tree)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func resetDB(db *pg.DB) error {
	opts := &orm.DropTableOptions{IfExists: true}

	for _, model := range []interface{}{(*User)(nil), (*Address)(nil), (*Tree)(nil)} {
		err := db.DropTable(model, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

// Initialize the database. Will reset the database if tables already exist.
func initDB() *pg.DB {
	db := connectDB()
	defer db.Close()

	err := resetDB(db)
	if err != nil {
		panic(err)
	}

	err = createSchema(db)
	if err != nil {
		panic(err)
	}

	user1 := &User{
		ID:       0,
		Username: "aleksey",
		Password: "1234",
		Admin:    true,
	}
	err = db.Insert(user1)
	if err != nil {
		panic(err)
	}

	return db
}

// Connect to the database. Note that this operation returns a live database connection, it is
// the user's responsibility to close this before exiting.
func connectDB() *pg.DB {
	opts := getConnectOpts()
	db := pg.Connect(opts)
	return db
}

//////////////
// DATA OPS //
//////////////

func createUser(conn *pg.DB, u User) error {
	err := conn.Insert(u)
	return err
}

// TODO: make use of this example code for useful things
// story1 := &Story{
// 	Title:    "Cool story",
// 	AuthorId: user1.Id,
// }
// err = db.Insert(story1)
// if err != nil {
// 	panic(err)
// }

// // Select user by primary key.
// user := &User{Id: user1.Id}
// err = db.Select(user)
// if err != nil {
// 	panic(err)
// }

// // Select all users.
// var users []User
// err = db.Model(&users).Select()
// if err != nil {
// 	panic(err)
// }

// // Select story and associated author in one query.
// story := new(Story)
// err = db.Model(story).
// 	Relation("Author").
// 	Where("story.id = ?", story1.Id).
// 	Select()
// if err != nil {
// 	panic(err)
// }

func main() {
	// initDB()
	conn := connectDB()
	// The following throws an error, for some reason?
	// admin := &User{Admin: true}
	admin := &User{ID: 1}
	err := conn.Select(admin)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", admin)
	conn.Close()
}

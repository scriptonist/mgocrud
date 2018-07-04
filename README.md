# `MgoCRUD`
generate mongo CRUD methods from  a golang struct

![Demo](https://raw.githubusercontent.com/scriptonist/mgocrud/master/artifacts/demo-mgocrud.gif)

# What is It ?
`mgocrud` is a cli tool which can generate basic CRUD functions from a `golang` struct
#### Example

Following are contents of `main.go`

```
package main

type User struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

```
Now on running `mgocrud generate main.go`, the file be appended with CRUD functions for the `User` struct.

```
package main

type User struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

func (u *User) Create(mgosession *mgo.Session, database, collection string) error {
	err := mgosession.DB(database).C(collection).Insert(u)
	if err != nil {
		return err
	}
	return err
}
func (u *User) Read(mgosession *mgo.Session, database, collection string, selector bson.M) (*[]User, error) {
	var results []User
	err := mgosession.DB(database).C(collection).Find(selector).All(&results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
func (u *User) Update(mgosession *mgo.Session, database, collection string, selector, change bson.M) error {
	err := mgosession.DB(database).C(collection).Update(selector, bson.M{"$set": change})
	if err != nil {
		return err
	}
	return nil
}
func (u *User) Delete(mgosession *mgo.Session, database, collection string, selector bson.M) (*mgo.ChangeInfo, error) {
	info, err := mgosession.DB(database).C(collection).RemoveAll(selector)
	if err != nil {
		return nil, err
	}
	return info, nil
}

```
`mgocrud` expects `github.com/globalsign/mgo` & `github.com/globalsign/bson` to be used as the driver for connecting to `mongo`.
# Okay, How can I try It ?

`go get github.com/scriptonist/mgocrud`

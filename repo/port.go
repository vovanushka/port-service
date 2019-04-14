package repo

import (
	"github.com/vovanushka/port-service/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PortRepo struct {
	collection *mgo.Collection
}

//Creates unque index on id field.
func idIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

// NewPortRepo is PortRepo constructor
func NewPortRepo(session *Session, dbName string, collectionName string) *PortRepo {
	// Connect to collection
	collection := session.GetCollection(dbName, collectionName)
	// Ensure that id index was created
	collection.EnsureIndex(idIndex())
	return &PortRepo{collection}
}

// Create function creates new port in the db. If the port
// with given id already exists, update it with new data.
func (r *PortRepo) Create(p *model.Port) error {
	_, err := r.collection.Upsert(bson.M{"id": p.ID}, bson.M{"$set": p})
	return err
}

// Get function return port from db by it's id.
// If there is no such id throws an error "not found"
func (r *PortRepo) Get(id string) (*model.Port, error) {
	port := &model.Port{}

	err := r.collection.Find(bson.M{"id": id}).One(port)
	if err != nil {
		return nil, err
	}

	return port, nil
}

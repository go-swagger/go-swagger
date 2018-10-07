package mgo_test

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	. "gopkg.in/check.v1"
)

type updateDesc struct {
	UpdatedFields map[string]interface{} `bson:"updatedFields"`
	RemovedFields []string               `bson:"removedFields"`
}

type evNamespace struct {
	DB   string `bson:"db"`
	Coll string `bson:"coll"`
}

type changeEvent struct {
	ID                interface{} `bson:"_id"`
	OperationType     string      `bson:"operationType"`
	FullDocument      *bson.Raw   `bson:"fullDocument,omitempty"`
	Ns                evNamespace `bson:"ns"`
	DocumentKey       M           `bson:"documentKey"`
	UpdateDescription *updateDesc `bson:"updateDescription,omitempty"`
}

func (s *S) TestStreamsWatch(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()
	coll := session.DB("mydb").C("mycoll")
	//add a mock document
	coll.Insert(M{"a": 0})

	pipeline := []bson.M{}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{})
	c.Assert(err, IsNil)

	err = changeStream.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsInsert(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add a mock document in order for the DB to be created
	err = coll.Insert(M{"a": 0})
	c.Assert(err, IsNil)

	//create the stream
	pipeline := []M{}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)

	//insert a new document
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 1})
	c.Assert(err, IsNil)
	//get the _id for later check
	type A struct {
		ID bson.ObjectId `bson:"_id"`
		A  int           `bson:"a"`
	}

	//get the event
	ev := changeEvent{}
	hasEvent := changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, true)

	//check event is correct
	oid := ev.DocumentKey["_id"].(bson.ObjectId)
	c.Assert(oid, Equals, id)
	c.Assert(ev.OperationType, Equals, "insert")
	c.Assert(ev.FullDocument, NotNil)
	a := A{}
	err = ev.FullDocument.Unmarshal(&a)
	c.Assert(err, IsNil)
	c.Assert(a.A, Equals, 1)
	c.Assert(ev.Ns.DB, Equals, "mydb")
	c.Assert(ev.Ns.Coll, Equals, "mycoll")

	err = changeStream.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsNextNoEventTimeout(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add a mock document in order for the DB to be created
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 0})
	c.Assert(err, IsNil)

	//create the stream
	pipeline := []M{}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)

	//check we timeout correctly on no events
	//we should get a false result and no error
	ev := changeEvent{}
	hasEvent := changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, false)
	c.Assert(changeStream.Err(), IsNil)
	c.Assert(changeStream.Timeout(), Equals, true)

	//test the same with default timeout (MaxTimeMS=1000)
	//create the stream
	changeStream, err = coll.Watch(pipeline, mgo.ChangeStreamOptions{})
	c.Assert(err, IsNil)
	hasEvent = changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, false)
	c.Assert(changeStream.Err(), IsNil)
	c.Assert(changeStream.Timeout(), Equals, true)

	err = changeStream.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsNextTimeout(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add a mock document in order for the DB to be created
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 0})
	c.Assert(err, IsNil)

	//create the stream
	pipeline := []M{}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)

	//insert a new document to trigger an event
	id = bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 1})
	c.Assert(err, IsNil)

	//ensure we get the event
	ev := changeEvent{}
	hasEvent := changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, true)

	//check we timeout correctly on no subsequent events
	//we should get a false result and no error
	ev = changeEvent{}
	hasEvent = changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, false)
	c.Assert(changeStream.Err(), IsNil)
	c.Assert(changeStream.Timeout(), Equals, true)

	//insert a new document to trigger an event
	id = bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 1})
	c.Assert(err, IsNil)

	//ensure we get the event
	ev = changeEvent{}
	hasEvent = changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, true)

	err = changeStream.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsDelete(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add a mock document in order for the DB to be created
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 0})
	c.Assert(err, IsNil)

	//create the changeStream
	pipeline := []M{}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)

	//delete the document
	err = coll.Remove(M{"_id": id})
	c.Assert(err, IsNil)

	//get the event
	ev := changeEvent{}
	hasEvent := changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, true)

	//check event is correct
	oid := ev.DocumentKey["_id"].(bson.ObjectId)
	c.Assert(oid, Equals, id)
	c.Assert(ev.OperationType, Equals, "delete")
	c.Assert(ev.FullDocument, IsNil)
	c.Assert(ev.Ns.DB, Equals, "mydb")
	c.Assert(ev.Ns.Coll, Equals, "mycoll")

	err = changeStream.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsUpdate(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add a mock document in order for the DB to be created
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 0, "toremove": 2})
	c.Assert(err, IsNil)

	//create the stream
	pipeline := []M{}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)

	//update document
	err = coll.UpdateId(id, M{"$set": M{"a": 1}, "$unset": M{"toremove": ""}})
	c.Assert(err, IsNil)

	//get the event
	ev := changeEvent{}
	hasEvent := changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, true)

	//check event is correct
	oid := ev.DocumentKey["_id"].(bson.ObjectId)
	c.Assert(oid, Equals, id)
	c.Assert(ev.OperationType, Equals, "update")
	c.Assert(ev.FullDocument, IsNil)
	c.Assert(len(ev.UpdateDescription.UpdatedFields), Equals, 1)
	c.Assert(len(ev.UpdateDescription.RemovedFields), Equals, 1)
	c.Assert(ev.UpdateDescription.UpdatedFields["a"], Equals, 1)
	c.Assert(ev.UpdateDescription.RemovedFields[0], Equals, "toremove")
	c.Assert(ev.Ns.DB, Equals, "mydb")
	c.Assert(ev.Ns.Coll, Equals, "mycoll")

	err = changeStream.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsUpdateFullDocument(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add a mock document in order for the DB to be created
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 0, "toremove": "bla"})
	c.Assert(err, IsNil)

	//create the stream
	pipeline := []M{}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500, FullDocument: mgo.UpdateLookup})
	c.Assert(err, IsNil)

	//update document
	err = coll.UpdateId(id, M{"$set": M{"a": 1}, "$unset": M{"toremove": ""}})
	c.Assert(err, IsNil)

	//get the event
	ev := changeEvent{}
	hasEvent := changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, true)

	type A struct {
		A        int     `bson:"a"`
		ToRemove *string `bson:"toremove"`
	}

	//check event is correct
	oid := ev.DocumentKey["_id"].(bson.ObjectId)
	c.Assert(oid, Equals, id)
	c.Assert(ev.OperationType, Equals, "update")
	c.Assert(len(ev.UpdateDescription.UpdatedFields), Equals, 1)
	c.Assert(len(ev.UpdateDescription.RemovedFields), Equals, 1)
	c.Assert(ev.UpdateDescription.UpdatedFields["a"], Equals, 1)
	c.Assert(ev.UpdateDescription.RemovedFields[0], Equals, "toremove")

	c.Assert(ev.FullDocument, NotNil)
	a := A{}
	err = ev.FullDocument.Unmarshal(&a)
	c.Assert(err, IsNil)
	c.Assert(a.A, Equals, 1)
	c.Assert(a.ToRemove, IsNil)
	c.Assert(ev.Ns.DB, Equals, "mydb")
	c.Assert(ev.Ns.Coll, Equals, "mycoll")

	err = changeStream.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsUpdateWithPipeline(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add two docs
	id1 := bson.NewObjectId()
	err = coll.Insert(M{"_id": id1, "a": 1})
	c.Assert(err, IsNil)
	id2 := bson.NewObjectId()
	err = coll.Insert(M{"_id": id2, "a": 2})
	c.Assert(err, IsNil)

	pipeline1 := []M{M{"$match": M{"documentKey._id": id1}}}
	changeStream1, err := coll.Watch(pipeline1, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)
	pipeline2 := []M{M{"$match": M{"documentKey._id": id2}}}
	changeStream2, err := coll.Watch(pipeline2, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)

	//update documents
	_, err = coll.UpdateAll(M{"_id": M{"$in": []bson.ObjectId{id1, id2}}}, M{"$inc": M{"a": 1}})
	c.Assert(err, IsNil)

	got1 := false
	got2 := false

	//check we got the update for id1 (and no other)
	for i := 0; i < 2; i++ {
		ev := changeEvent{}
		hasEvent := changeStream1.Next(&ev)
		//we will accept only one event, the one that corresponds to our id1
		c.Assert(got1 && hasEvent, Equals, false)
		if hasEvent {
			oid := ev.DocumentKey["_id"].(bson.ObjectId)
			c.Assert(oid, Equals, id1)
			got1 = true
		}
	}
	c.Assert(got1, Equals, true)

	//check we got the update for id2 (and no other)
	for i := 0; i < 2; i++ {
		ev := changeEvent{}
		hasEvent := changeStream2.Next(&ev)
		//we will accept only one event, the one that corresponds to our id2
		c.Assert(got2 && hasEvent, Equals, false)
		if hasEvent {
			oid := ev.DocumentKey["_id"].(bson.ObjectId)
			c.Assert(oid, Equals, id2)
			got2 = true
		}
	}
	c.Assert(got2, Equals, true)

	err = changeStream1.Close()
	c.Assert(err, IsNil)
	err = changeStream2.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsResumeTokenMissingError(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add a mock document in order for the DB to be created
	err = coll.Insert(M{"a": 0})
	c.Assert(err, IsNil)

	//create the stream
	pipeline := []M{{"$project": M{"_id": 0}}}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)

	//insert a new document
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 1})
	c.Assert(err, IsNil)

	//check we get the correct error
	ev := changeEvent{}
	hasEvent := changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, false)
	c.Assert(changeStream.Err().Error(), Equals, "resume token missing from result")

	err = changeStream.Close()
	c.Assert(err, IsNil)
}

func (s *S) TestStreamsClosedStreamError(c *C) {
	if !s.versionAtLeast(3, 6) {
		c.Skip("ChangeStreams only work on 3.6+")
	}
	session, err := mgo.Dial("localhost:40011")
	c.Assert(err, IsNil)
	defer session.Close()

	coll := session.DB("mydb").C("mycoll")

	//add a mock document in order for the DB to be created
	err = coll.Insert(M{"a": 0})
	c.Assert(err, IsNil)

	//create the stream
	pipeline := []M{{"$project": M{"_id": 0}}}
	changeStream, err := coll.Watch(pipeline, mgo.ChangeStreamOptions{MaxAwaitTimeMS: 1500})
	c.Assert(err, IsNil)

	//insert a new document
	id := bson.NewObjectId()
	err = coll.Insert(M{"_id": id, "a": 1})
	c.Assert(err, IsNil)

	err = changeStream.Close()
	c.Assert(err, IsNil)

	//check we get the correct error
	ev := changeEvent{}
	hasEvent := changeStream.Next(&ev)
	c.Assert(hasEvent, Equals, false)
	c.Assert(changeStream.Err().Error(), Equals, "illegal use of a closed ChangeStream")
}

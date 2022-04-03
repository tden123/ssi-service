package storage

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBoltDB(t *testing.T) {
	db, err := NewBoltDBWithFile("test.db")
	assert.NoError(t, err)
	assert.NotEmpty(t, db)

	t.Cleanup(func() {
		_ = db.Close()
		os.Remove("test.db")
	})

	// create a name space and a message in it
	namespace := "F1"

	team1 := "Red Bull"
	players1 := []string{"Max Verstappen", "Sergio Pérez"}
	p1Bytes, err := json.Marshal(players1)
	assert.NoError(t, err)

	err = db.Write(namespace, team1, p1Bytes)
	assert.NoError(t, err)

	// get it back
	gotPlayers1, err := db.Read(namespace, team1)
	assert.NoError(t, err)

	var players1Result []string
	err = json.Unmarshal(gotPlayers1, &players1Result)
	assert.NoError(t, err)
	assert.EqualValues(t, players1, players1Result)

	// get a value from a namespace that doesn't exist
	_, err = db.Read("bad", "worse")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace<bad> does not exist")

	// get a value that doesn't exist in the namespace
	noValue, err := db.Read(namespace, "Porsche")
	assert.NoError(t, err)
	assert.Empty(t, noValue)

	// create a second value in the namespace
	team2 := "McLaren"
	players2 := []string{"Lando Norris", "Daniel Ricciardo"}
	p2Bytes, err := json.Marshal(players2)
	assert.NoError(t, err)

	err = db.Write(namespace, team2, p2Bytes)
	assert.NoError(t, err)

	// get all values from the namespace
	gotAll, err := db.ReadAll(namespace)
	assert.NoError(t, err)
	assert.True(t, len(gotAll) == 2)

	_, gotRedBull := gotAll[team1]
	assert.True(t, gotRedBull)

	_, gotMcLaren := gotAll[team2]
	assert.True(t, gotMcLaren)

	// delete value in the namespace
	err = db.Delete(namespace, team2)
	assert.NoError(t, err)

	gotPlayers2, err := db.Read(namespace, team2)
	assert.NoError(t, err)
	assert.Empty(t, gotPlayers2)

	// delete value in a namespace that doesn't exist
	err = db.Delete("bad", team2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace<bad> does not exist")

	// delete a namespace that doesn't exist
	err = db.DeleteNamespace("bad")
	assert.Contains(t, err.Error(), "could not delete namespace<bad>")

	// delete namespace
	err = db.DeleteNamespace(namespace)
	assert.NoError(t, err)

	_, err = db.Read(namespace, team1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace<F1> does not exist")
}

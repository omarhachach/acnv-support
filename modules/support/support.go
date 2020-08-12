package support

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/rs/xid"
)

// Entry is a support message.
type Entry struct {
	ID             string `json:"-"`
	MessageContent string `json:"message_content"`
	ChannelID      string `json:"channel_id"`
	SenderID       string `json:"sender_id"`
}

var ErrNotFound = errors.New("the entry wasn't found")

type Map map[string]Entry

func (s Map) Insert(entry Entry) Entry {
	guid := xid.New().String()

	_, ok := s[guid]
	if ok {
		return s.Insert(entry)
	}

	entry.ID = guid
	s[guid] = entry

	return entry
}

func (m *Module) AddEntryToSupportFile(entry Entry) (Entry, error) {
	m.Lock()
	defer m.Unlock()

	supportMap, err := m.getSupportMapFromFile()
	if err != nil {
		return Entry{}, err
	}

	newEntry := supportMap.Insert(entry)

	supportBytes, err := json.Marshal(supportMap)
	if err != nil {
		return Entry{}, fmt.Errorf("error writing to support map file %w", err)
	}

	err = ioutil.WriteFile("support-file.json", supportBytes, 0755)
	if err != nil {
		return Entry{}, err
	}

	return newEntry, nil
}

func (m *Module) getSupportMapFromFile() (Map, error) {
	fileBytes, err := ioutil.ReadFile("support-file.json")
	if err != nil {
		return nil, err
	}

	var supportMap Map

	err = json.Unmarshal(fileBytes, &supportMap)
	if err != nil {
		return nil, err
	}

	return supportMap, nil
}

func (m *Module) GetEntryFromSupportFile(entryId string) (Entry, error) {
	supportMap, err := m.getSupportMapFromFile()
	if err != nil {
		return Entry{}, err
	}

	entry, ok := supportMap[entryId]
	if !ok {
		return Entry{}, ErrNotFound
	}

	entry.ID = entryId

	return entry, nil
}

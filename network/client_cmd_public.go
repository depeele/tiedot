/* tiedot command implementations - client side, public APIs. */
package network

import (
	"encoding/json"
	"fmt"
)

// Tell server to shutdown, and shutdown myself (client) as well.
func (tc *Client) ShutdownServer() {
	tc.writeReq(SHUTDOWN)
	tc.ShutdownClient()
}

// Create a collection.
func (tc *Client) ColCreate(name string, numParts int) error {
	return tc.getOK(fmt.Sprintf("%s %s %d", COL_CREATE, name, numParts))
}

// Get all collection information (collection name VS number of partitions).
func (tc *Client) ColAll() (all map[string]int, err error) {
	all = make(map[string]int)
	js, err := tc.getJSON(COL_ALL)
	if err != nil {
		return
	}
	for name, numParts := range js.(map[string]interface{}) {
		all[name] = int(numParts.(float64))
	}
	return
}

// Rename a collection.
func (tc *Client) ColRename(oldName, newName string) error {
	return tc.getOK(fmt.Sprintf("%s %s %s", COL_RENAME, oldName, newName))
}

// Drop a collection.
func (tc *Client) ColDrop(colName string) error {
	return tc.getOK(fmt.Sprintf("%s %s", COL_DROP, colName))
}

// Create an index.
func (tc *Client) IdxCreate(colName, idxPath string) error {
	return tc.getOK(fmt.Sprintf("%s %s %s", IDX_CREATE, colName, idxPath))
}

// Get all indexed paths.
func (tc *Client) IdxAll(colName string) (paths []string, err error) {
	js, err := tc.getJSON(fmt.Sprintf("%s %s", IDX_ALL, colName))
	if err != nil {
		return
	}
	paths = make([]string, 0, 12)
	for _, path := range js.([]interface{}) {
		paths = append(paths, path.(string))
	}
	return paths, nil
}

// Drop an index.
func (tc *Client) IdxDrop(colName, idxPath string) error {
	return tc.getOK(fmt.Sprintf("%s %s %s", IDX_DROP, colName, idxPath))
}

// Insert a document, return its ID.
func (tc *Client) ColInsert(colName string, js map[string]interface{}) (uint64, error) {
	if serialized, err := json.Marshal(js); err != nil {
		return 0, err
	} else {
		return tc.getUint64(fmt.Sprintf("%s %s %s", COL_INSERT, colName, string(serialized)))
	}
}

// Get a document by ID.
func (tc *Client) ColGet(colName string, id uint64) (doc interface{}, err error) {
	return tc.getJSON(fmt.Sprintf("%s %s %d", COL_GET, colName, id))
}

// Get a document by ID.
func (tc *Client) ColGetJS(colName string, id uint64) (doc string, err error) {
	return tc.getStr(fmt.Sprintf("%s %s %d", COL_GET, colName, id))
}

// Update a document by ID.
func (tc *Client) ColUpdate(colName string, id uint64, js map[string]interface{}) (err error) {
	if serialized, err := json.Marshal(js); err != nil {
		return err
	} else {
		return tc.getOK(fmt.Sprintf("%s %s %d %s", COL_UPDATE, colName, id, string(serialized)))
	}
}

// Delete a document by ID.
func (tc *Client) ColDelete(colName string, id uint64) (err error) {
	return tc.getOK(fmt.Sprintf("%s %s %d", COL_DELETE, colName, id))
}
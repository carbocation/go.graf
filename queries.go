package main

var queries = struct {
	FindOneUserById string
	FindOneByHandle	string
}{
	FindOneUserById: `SELECT id, handle, password, created FROM account WHERE id=$1 LIMIT 1`,
	FindOneByHandle: `SELECT id, handle, password, created FROM account WHERE handle=$1 LIMIT 1`, 
}

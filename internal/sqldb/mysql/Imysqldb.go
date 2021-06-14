package mysql

import "gorm.io/gorm"

type IMySqlConnector interface {
	GetConn() *gorm.DB
	Ping() error
}

var MySqlConnector IMySqlConnector
var MySqlGoAutodialConnector IMySqlConnector

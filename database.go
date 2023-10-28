package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// go...
//          _____                _____                    _____                   _______                   _____
//         /\    \              |\    \                  /\    \                 /::\    \                 /\    \
//        /::\____\             |:\____\                /::\    \               /::::\    \               /::\____\
//       /::::|   |             |::|   |               /::::\    \             /::::::\    \             /:::/    /
//      /:::::|   |             |::|   |              /::::::\    \           /::::::::\    \           /:::/    /
//     /::::::|   |             |::|   |             /:::/\:::\    \         /:::/~~\:::\    \         /:::/    /
//    /:::/|::|   |             |::|   |            /:::/__\:::\    \       /:::/    \:::\    \       /:::/    /
//   /:::/ |::|   |             |::|   |            \:::\   \:::\    \     /:::/    / \:::\    \     /:::/    /
//  /:::/  |::|___|______       |::|___|______    ___\:::\   \:::\    \   /:::/____/   \:::\____\   /:::/    /
// /:::/   |::::::::\    \      /::::::::\    \  /\   \:::\   \:::\    \ |:::|    |     |:::|    | /:::/    /
///:::/    |:::::::::\____\    /::::::::::\____\/::\   \:::\   \:::\____\|:::|____|     |:::|____|/:::/____/
//\::/    / ~~~~~/:::/    /   /:::/~~~~/~~      \:::\   \:::\   \::/    / \:::\   _\___/:::/    / \:::\    \
// \/____/      /:::/    /   /:::/    /          \:::\   \:::\   \/____/   \:::\ |::| /:::/    /   \:::\    \
//             /:::/    /   /:::/    /            \:::\   \:::\    \        \:::\|::|/:::/    /     \:::\    \
//            /:::/    /   /:::/    /              \:::\   \:::\____\        \::::::::::/    /       \:::\    \
//           /:::/    /    \::/    /                \:::\  /:::/    /         \::::::::/    /         \:::\    \
//          /:::/    /      \/____/                  \:::\/:::/    /           \::::::/    /           \:::\    \
//         /:::/    /                                 \::::::/    /             \::::/____/             \:::\    \
//        /:::/    /                                   \::::/    /               |::|    |               \:::\____\
//        \::/    /                                     \::/    /                |::|____|                \::/    /
//         \/____/                                       \/____/                  ~~                       \/____/

type MysqlDb struct {
	dsn    string
	taber  *sql.DB
	basket *sql.DB
}

var mysql MysqlDb

// go go go
func (_mysql MysqlDb) Init() MysqlDb {
	mysql.dsn = env("DSN_MYSQL", "#_PRIVATE_#")
	mysql.taber = mysql.Connect(mysql.dsn, "#_PRIVATE_#")
	mysql.basket = mysql.Connect(mysql.dsn, "#_PRIVATE_#")
	return mysql
}

// go connect
func (_mysql MysqlDb) Connect(dsn string, database string) *sql.DB {
	db, err := sql.Open("mysql", dsn+"/"+database)
	if err != nil || db == nil {
		panic("Ошибка подключения к MySQL: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic("Ошибка подключения к MySQL: " + err.Error())
	}

	db.SetConnMaxLifetime(time.Second * 30)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.Query("SET NAMES `utf8`")

	if database == "taber" {
		mysql.taber = db
	}
	if database == "basket" {
		mysql.basket = db
	}
	return db
}

// go use
func (_mysql MysqlDb) use(database string) *sql.DB {
	if database == "taber" {
		if mysql.taber != nil {
			return mysql.taber
		}
		return mysql.Init().taber
	}
	if database == "basket" {
		if mysql.basket != nil {
			return mysql.basket
		}
		return mysql.Init().basket
	}
	return nil
}

// go out
func (_mysql MysqlDb) Close() {
	mysql.taber.Close()
	mysql.basket.Close()
}

// go ...
//           _____                   _______                   _____                    _____                   _______
//          /\    \                 /::\    \                 /\    \                  /\    \                 /::\    \
//         /::\____\               /::::\    \               /::\____\                /::\    \               /::::\    \
//        /::::|   |              /::::::\    \             /::::|   |               /::::\    \             /::::::\    \
//       /:::::|   |             /::::::::\    \           /:::::|   |              /::::::\    \           /::::::::\    \
//      /::::::|   |            /:::/~~\:::\    \         /::::::|   |             /:::/\:::\    \         /:::/~~\:::\    \
//     /:::/|::|   |           /:::/    \:::\    \       /:::/|::|   |            /:::/  \:::\    \       /:::/    \:::\    \
//    /:::/ |::|   |          /:::/    / \:::\    \     /:::/ |::|   |           /:::/    \:::\    \     /:::/    / \:::\    \
//   /:::/  |::|___|______   /:::/____/   \:::\____\   /:::/  |::|   | _____    /:::/    / \:::\    \   /:::/____/   \:::\____\
//  /:::/   |::::::::\    \ |:::|    |     |:::|    | /:::/   |::|   |/\    \  /:::/    /   \:::\ ___\ |:::|    |     |:::|    |
// /:::/    |:::::::::\____\|:::|____|     |:::|    |/:: /    |::|   /::\____\/:::/____/  ___\:::|    ||:::|____|     |:::|    |
// \::/    / ~~~~~/:::/    / \:::\    \   /:::/    / \::/    /|::|  /:::/    /\:::\    \ /\  /:::|____| \:::\    \   /:::/    /
//  \/____/      /:::/    /   \:::\    \ /:::/    /   \/____/ |::| /:::/    /  \:::\    /::\ \::/    /   \:::\    \ /:::/    /
//              /:::/    /     \:::\    /:::/    /            |::|/:::/    /    \:::\   \:::\ \/____/     \:::\    /:::/    /
//             /:::/    /       \:::\__/:::/    /             |::::::/    /      \:::\   \:::\____\        \:::\__/:::/    /
//            /:::/    /         \::::::::/    /              |:::::/    /        \:::\  /:::/    /         \::::::::/    /
//           /:::/    /           \::::::/    /               |::::/    /          \:::\/:::/    /           \::::::/    /
//          /:::/    /             \::::/    /                /:::/    /            \::::::/    /             \::::/    /
//         /:::/    /               \::/____/                /:::/    /              \::::/    /               \::/____/
//         \::/    /                 ~~                      \::/    /                \::/____/                 ~~
//          \/____/                                           \/____/

type MongoDb struct {
	dsn        string
	context    context.Context
	connection *mongo.Client
	taber      *mongo.Database
}

var mongodb MongoDb

// go go go
func (_mongodb MongoDb) Init() MongoDb {
	mongodb.dsn = env("DSN_MONGO", "#_PRIVATE_#")
	mongodb.context = context.Background()
	mongodb.taber = mongodb.Connect(mongodb.dsn, "#_PRIVATE_#")
	return mongodb
}

// go connect
func (_mongodb MongoDb) Connect(dsn string, database string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(dsn)
	mongodb.connection, _ = mongo.Connect(mongodb.context, clientOptions)

	err := mongodb.connection.Ping(mongodb.context, nil)
	if err != nil {
		panic(err.Error())
	}

	mongodb.taber = mongodb.connection.Database("taber")
	return mongodb.taber
}

// go get collection
func (_mongodb MongoDb) collection(name string, new ...bool) *mongo.Collection {
	if mongodb.taber == nil {
		mongodb.Init()
	}
	if new != nil && new[0] {
		existCollection := mongodb.taber.Collection(name)
		if existCollection != nil {
			existCollection.Drop(mongodb.context)
		}
		mongodb.taber.CreateCollection(mongodb.context, name)
	}
	return mongodb.taber.Collection(name)
}

// go get index
func (_mongodb MongoDb) index(collection string, index string) *mongo.Collection {
	if mongodb.taber == nil {
		mongodb.Init()
	}
	table := mongodb.taber.Collection(collection)
	if table == nil {
		panic("Не существует коллекции с именем " + collection)
	}
	table.Indexes().CreateOne(
		mongodb.context,
		mongo.IndexModel{Keys: bson.D{{Key: index, Value: 1}}},
	)
	return table
}

// go out
func (_mongodb MongoDb) Close() {
	mongodb.taber.Client().Disconnect(mongodb.context)
}

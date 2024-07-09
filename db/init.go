package db

func InitDB() {
	ConnectDB()
	InitQueryBuilder()
	InitRedis()
	InitReaderRepo()
	InitBookRepo()
	InitBookReqeustRepo()
	InitAdminRepo()
}

package gocrud_redis

import (
	"github.com/kordar/gocrud"
	"sync"
)

var initOnce sync.Once

func InitExec() {
	initOnce.Do(func() {
		registerExec("redis")
		registerExec("goredis")
	})
}

func registerExec(driver string) {
	gocrud.AddExecute("=", EQ, driver)
	gocrud.AddExecute("EQ", EQ, driver)
	gocrud.AddExecute("!=", NEQ, driver)
	gocrud.AddExecute("<>", NEQ, driver)
	gocrud.AddExecute("NEQ", NEQ, driver)
	gocrud.AddExecute("LT", LT, driver)
	gocrud.AddExecute("<", LT, driver)
	gocrud.AddExecute("LE", LE, driver)
	gocrud.AddExecute("<=", LE, driver)
	gocrud.AddExecute("GT", GT, driver)
	gocrud.AddExecute(">", GT, driver)
	gocrud.AddExecute("GE", GE, driver)
	gocrud.AddExecute(">=", GE, driver)
	gocrud.AddExecute("IN", IN, driver)
	gocrud.AddExecute("NOTIN", NOTIN, driver)
	gocrud.AddExecute("LIKE", LIKE, driver)
	gocrud.AddExecute("NOTLIKE", NOTLIKE, driver)
	gocrud.AddExecute("LIKELEFT", LIKELEFT, driver)
	gocrud.AddExecute("LIKERIGHT", LIKERIGHT, driver)
	gocrud.AddExecute("BETWEEN", BETWEEN, driver)
	gocrud.AddExecute("NOTBETWEEN", NOTBETWEEN, driver)
	gocrud.AddExecute("ISNULL", ISNULL, driver)
	gocrud.AddExecute("ISNOTNULL", ISNOTNULL, driver)

	gocrud.AddExecute("ASC", ASC, driver)
	gocrud.AddExecute("DESC", DESC, driver)

	gocrud.AddExecute("SAVE", SAVE, driver)
	gocrud.AddExecute("SETVAL", SETVAL, driver)
	gocrud.AddExecute("UPDATE", UPDATE, driver)
	gocrud.AddExecute("UPDATES", UPDATES, driver)
	gocrud.AddExecute("CREATE", CREATE, driver)
	gocrud.AddExecute("PAGE", PAGE, driver)
	gocrud.AddExecute("DELETE", DELETE, driver)
}

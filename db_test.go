package dbdao

import (
	"testing"

	"github.com/hy0kl/gconfig"
)

func TestDb(t *testing.T) {
	gconfig.SetConfigFile("./conf/conf.ini")

	db := GetDbInstance("blog")

	masterEngine := db.Engine.Master()
	results, err := masterEngine.Query("select 123 as name")
	if err != nil {
		t.Error("fail", err)
	}

	if len(results) == 1 && string(results[0]["name"]) == "123" {
		t.Log("pass")
	} else {
		t.Error("fail")
	}
}

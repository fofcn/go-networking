package gin_test

import (
	"os"
	"testing"
)

func TestLoadDBConfig(t *testing.T) {
	// 设置环境变量
	os.Setenv("DB_DRIVER", "mysql_test")
	os.Setenv("DB_HOST", "localhost_test")
	os.Setenv("DB_USER", "root_test")
	os.Setenv("DB_PASSWORD", "password_test")
	os.Setenv("DB_NAME", "mydb_test")
	os.Setenv("DB_PORT", "3306_test")

	// 调用函数获取返回值
	config := LoadDBConfig()

	// 断言返回值与预期值是否相等
	if config.DbDriver != "mysql_test" {
		t.Errorf("DbDriver expected: mysql_test, got: %s", config.DbDriver)
	}
	if config.DbHost != "localhost_test" {
		t.Errorf("DbHost expected: localhost_test, got: %s", config.DbHost)
	}
	if config.DbUser != "root_test" {
		t.Errorf("DbUser expected: root_test, got: %s", config.DbUser)
	}
	if config.DbPassword != "password_test" {
		t.Errorf("DbPassword expected: password_test, got: %s", config.DbPassword)
	}
	if config.DbName != "mydb_test" {
		t.Errorf("DbName expected: mydb_test, got: %s", config.DbName)
	}
	if config.DbPort != "3306_test" {
		t.Errorf("DbPort expected: 3306_test, got: %s", config.DbPort)
	}
}

func TestLoadDBConfig_DefaultValues(t *testing.T) {
	// 清除环境变量
	os.Unsetenv("DB_DRIVER")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PORT")

	// 调用函数获取返回值
	config := LoadDBConfig()

	// 断言返回值与默认值是否相等
	if config.DbDriver != "mysql" {
		t.Errorf("DbDriver expected: mysql, got: %s", config.DbDriver)
	}
	if config.DbHost != "localhost" {
		t.Errorf("DbHost expected: localhost, got: %s", config.DbHost)
	}
	if config.DbUser != "root" {
		t.Errorf("DbUser expected: root, got: %s", config.DbUser)
	}
	if config.DbPassword != "password" {
		t.Errorf("DbPassword expected: password, got: %s", config.DbPassword)
	}
	if config.DbName != "mydb" {
		t.Errorf("DbName expected: mydb, got: %s", config.DbName)
	}
	if config.DbPort != "3306" {
		t.Errorf("DbPort expected: 3306, got: %s", config.DbPort)
	}
}

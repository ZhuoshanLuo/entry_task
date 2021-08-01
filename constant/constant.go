package constant

//常量
const (
	Driver = "mysql"
	SqlUser  = "root"
	Passwd   = "Luo2566288"
	Host     = "127.0.0.1:3306"
	Database = "et_db"
	Salt     = "abcd"
	PasswdPattern = "^(?=.*[A-Za-z])(?=.*\\d)[A-Za-z\\d]{8,}$" // 8位密码，至少包含一个字母和数字
	NamePattern = "^[^0-9][\\w a-z A-Z 0-9]" //不能是数字开头，可包含数字大小写英文
)


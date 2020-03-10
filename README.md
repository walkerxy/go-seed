# go-seed
go语言版的数据填充，目前只支持mysql数据库及yaml文件

# Install
go get github.com/walkerxy/go-seed

# Quick Start
- 在目录下新建一个conf.yaml配置数据库信息
``` 
database:
    host: "192.168.1.1"
    port: "3306"
    database_name: ""
    username: root
    password: 
    table_prefix: ""
```
- 基本使用
```
package main

func main() {
    filepath := "./seeds"
    filename := "user.yaml"
    seed := NewSeed(filepath, filename)
    seed.SetTablePrefix("mc")
    seed.Fill()
}
```

- 多文件
```
package main

func seedFile(filepath string, filename string) {
	seed := NewSeed(filepath, filename)
	seed.SetTablePrefix("mc")
	seed.Fill()
}

func main() {
    dir := "./seeds"
    
    fileChan := make(chan string)
	go func() {
		WalkDir(dir, "yaml", fileChan)
		close(fileChan)
	}()

	for file := range fileChan {
		seedFile(dir, file)
	}
}
```

# License
MIT

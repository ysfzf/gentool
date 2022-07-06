#### Give a star before you see it. Ha ha ha ~ ~

Generates a protobuf file from your mysql database,May also generates gorm files from your database.

### protobuf
 
 ```
 gentool proto --config config.yaml  

 ```

 config.yaml

```yaml
dbType: mysql   
host: localhost 
port: 3306
user: root
password: 123456
schema: ff
tables: 
  - users
serviceName: user
packageName: user
goPackageName: ./user
ignoreTables:  
  - user_logs
ignoreColumns: 
  - id
  - created_at
  - updated_at
  - deleted_at
  - version

```

### gorm 

```
 gentool gorm --config config.yaml 
```

config.yaml

```yaml
dbType: mysql
host: localhost
port: 3306
user: root
password: 123456
schema: fcar
tables: 
  - users
  - merchant_users
onlyModel: false
outPath: ./dao/query  #Cannot end with "model"
outFile:
withUnitTest: false
modelPkgName: model
fieldNullable: false
fieldWithIndexTag: false
fieldWithTypeTag: false

```

#### Thanks 
    go-gorm : https://github.com/go-gorm/gen
    Mikaelemmmm : https://github.com/zeromicro/go-zero

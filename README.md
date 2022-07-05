#### Give a star before you see it. Ha ha ha ~ ~

Generates a protobuf file from your mysql database.

### Uses
 
 ```
 sql2pb  ./config.yaml >test.proto

 ```

### yaml

```yaml
dbType: mysql   
host: localhost 
port: 3306
user: root
password: 123456
schema: ff
tables: 
serviceName: merchant
packageName: merchant
goPackageName: ./merchant
ignoreTables:  
  - user_logs
ignoreColumns: 
  - id
  - created_at
  - updated_at
  - deleted_at
  - version

```

#### Thanks for schemabuf
    schemabuf : https://github.com/mcos/schemabuf

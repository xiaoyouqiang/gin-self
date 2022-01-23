## 命令说明
根目录下执行：`./scripts/gen_model.sh addr user pass name tables`；
- addr：数据库地址，例如：127.0.0.1:3306
- user：账号，如：root
- pass：密码，如：root
- name：数据库名称，如：test
- tables：表名，* 表示所有表，多个表名可用 “,”分割，例如：test,test1

如：
```
./scripts/gen_model.sh 127.0.0.1:3306 root root test test1
```
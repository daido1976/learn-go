```sh
$ curl -X POST http://localhost:8080/todos/ -d '{ "id": 1, "title": "Test", "body": "bodyだよ" }' | jq .
```

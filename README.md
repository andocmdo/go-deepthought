# go-deepthought
RESTful JSON (and other formats) API Backend Job Server for my personal Deepthought cluster

# Example curl commands (these are outdated... will fix soon):
```

curl -H "Content-Type: application/json" -d '{"args":{"command":"sleep 5; echo hello","endDate":"2017-12-27","maxGen":"1000","mutRate":"0.04","popSize":"235","startDate":"2016-12-26","symbol":"AAPL"}}' http://localhost:8080/api/v1/jobs
curl -H "Content-Type: application/json" -d '{"symbol":"192.168.1.45","port":"12365"}' http://localhost:8080/api/v1/workers

```

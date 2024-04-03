## Simple url shortener

`go run ./cmd/` - Без использования бд

`go run ./cmd/ -d` - С использованием бд


**curl -X POST -d "url=http://cjdr17afeihmk.biz/123/kdni9/z9d112423421" http://localhost:8080**

response: http://localhost:8080/a0C1S

**curl localhost:8080/a0C1S**

response: http://cjdr17afeihmk.biz/123/kdni9/z9d11242342

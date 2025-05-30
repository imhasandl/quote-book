insert Quote

curl -X POST http://localhost:8080/quotes \
-H "Content-Type: application/json" \
-d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'

Get quotes

curl http://localhost:8080/quotes

Get random quote 

curl http://localhost:8080/quotes/random


Get quotes from author

curl http://localhost:8080/quotes?author=Confucius


Delete quote

curl -X DELETE http://localhost:8080/quotes/1
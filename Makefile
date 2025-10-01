all:

build:
	docker-compose build

up:
	docker-compose up

down:
	docker-compose down

logs:
	docker-compose logs -f sms-service user-service email-service analytics-service

request:
	curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"email": "test@example.com", "phone": "+1234567890"}'
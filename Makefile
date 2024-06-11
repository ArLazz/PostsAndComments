.PHONY: start tests

start: 
	docker-compose -f docker-compose.yml up --build

tests: 
	docker-compose -f tests/test_docker-compose.yml up --build
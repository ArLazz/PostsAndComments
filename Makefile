start: 
	docker-compose -f docker-compose.yml up --build

tests: 
	docker-compose -f test_docker-compose.yml up --build
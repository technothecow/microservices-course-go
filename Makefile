.PHONY: run gen

# Run the services
run:
	docker-compose -f docker-compose.yml -f gateway/docker-compose.yml -f users/docker-compose.yml -f posts/docker-compose.yml up --build

gen:
	python3 -m make-gen
run:
	go run .

build:
	docker build -t forum .

docker-run:
	docker run --network=host -p 4000:4000 --name cinemaForum forum 

stop:
	docker stop cinemaForum
	
delete:
	docker rm $$(docker ps -aq)
	docker system prune -a	
build:
	go build -o traymenu cmd/menu/menu.go

create-default-config:
	cp -r resources ~/.traymenu/
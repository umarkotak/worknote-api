run:
	go run .

migrate-up:
	go run . migrate up

bin:
	go build -o worknote-api .

bin_run:
	./worknote-api

nohup_run:
	nohup ./worknote-api &

stopd:
	pkill worknote-api

statusd:
	ps aux | grep worknote-api

logs:
	tail -f worknote-api.error.log

install-service:
	sudo chmod +x worknote-api
	sudo cp com.worknote-api.plist /Library/LaunchDaemons
	sudo chmod +x /Library/LaunchDaemons/com.worknote-api.plist
	sudo launchctl bootstrap system /Library/LaunchDaemons/com.worknote-api.plist

uninstall-service:
	sudo launchctl unload /Library/LaunchDaemons/com.worknote-api.plist
	sudo rm /Library/LaunchDaemons/com.worknote-api.plist

start:
	sudo launchctl start com.worknote-api

stop:
	sudo launchctl stop com.worknote-api

deploy:
	git pull --rebase origin master
	go mod tidy
	go mod vendor
	make bin
	make stop
	make start

status:
	sudo lsof -i :6020

export GO111MODULE=on
DB_HOST:=db
DB_PORT:=3306
DB_USER:=root
DB_PASS:=pass
DB_NAME:=isucon

MYSQL_CMD:=mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) -p$(DB_PASS) $(DB_NAME)

NGX_LOG:=/var/log/nginx/access.log
MYSQL_LOG:=/var/log/mysql/mysql_slow.log

KATARU_CFG:=./kataribe.toml

SLACKCAT:=slackcat --tee 
SLACKRAW:=slackcat 

PPROF:=go tool pprof -png -output pprof.png http://localhost:8080/debug/pprof/profile

PROJECT_ROOT:=/go/src/github.com/gosagawa/isucon
BUILD_DIR:=/go/src/github.com/gosagawa/isucon
BIN_NAME:=bin/isucon

CA:=-o /dev/null -s -w "%{http_code}\n"

all: build

# デバッグおよび負荷試験用。実際の競技時には不要。
# ここから -------------------------------------------------

.PHONY: ssh
ssh:
	docker-compose exec web bash

sshdb:
	docker-compose exec db bash

sshproxy:
	docker-compose exec proxy bash

.PHONY: st
st:
	echo "GET http://localhost:8080/user/index" | vegeta attack -rate=300 -duration=5s | tee vegeta/result/result.bin

stresult:
	vegeta report vegeta/result/result.bin | $(SLACKCAT)

# ここまで -------------------------------------------------

.PHONY: clean
clean:
	cd $(BUILD_DIR); \
	rm -rf torb

deps:
	cd $(BUILD_DIR); \
	go mod download

.PHONY: build
build:
	cd $(BUILD_DIR); \
	go build -o $(BIN_NAME)
	#TODO ビルドコマンドを正しいものに直す

.PHONY: restart
restart:

	#TODO リスタートコマンドがあれば記載

.PHONY: test
test:
	curl localhost $(CA)

.PHONY: dev
dev: build 
	cd $(BUILD_DIR); \
	./$(BIN_NAME)

.PHONY: bench-dev
bench-dev: commit before slow-on dev

.PHONY: bench
bench: commit before build restart log

.PHONY: log
log: 
	#TODO コマンド要修正
	#sudo journalctl -u isucari.golang -n10 -f

.PHONY: maji
bench: commit before build restart

.PHONY: anal
anal: slow kataru

.PHONY: push
push: 
	git push

.PHONY: commit
commit:
	cd $(PROJECT_ROOT); \
	git add .; \
	git commit --allow-empty -m "bench"

.PHONY: before
before:
	$(eval when := $(shell date "+%s"))
	mkdir -p ~/logs/$(when)
	@if [ -f $(NGX_LOG) ]; then \
		sudo mv -f $(NGX_LOG) ~/logs/$(when)/ ; \
	fi
	# @if [ -f $(MYSQL_LOG) ]; then \
	# 	sudo mv -f $(MYSQL_LOG) ~/logs/$(when)/ ; \
	# fi
	sudo systemctl restart nginx
	# sudo systemctl restart mysql

.PHONY: slow
slow: 
	sudo pt-query-digest $(MYSQL_LOG) | $(SLACKCAT)

.PHONY: kataru
kataru:
	sudo cat $(NGX_LOG) | kataribe -f ./kataribe.toml | $(SLACKCAT)

.PHONY: pprof
pprof:
	$(PPROF)
	$(SLACKRAW) -n pprof.png ./pprof.png
	mv pprof.png pprof/

.PHONY: mysql
mysql:
	sudo $(MYSQL_CMD) 

.PHONY: slow-on
slow-on:
	sudo $(MYSQL_CMD) -e "set global slow_query_log_file = '$(MYSQL_LOG)'; set global long_query_time = 0; set global slow_query_log = ON;"

.PHONY: slow-off
slow-off:
	sudo $(MYSQL_CMD) -e "set global slow_query_log = OFF;"

.PHONY: goosecreate
goosecreate:
	goose create mod sql

gooseup:
	goose up

.PHONY: goosedown
goosedown:
	goose down

.PHONY: setup
setup:
	sudo apt install -y percona-toolkit dstat git unzip snapd
	mkdir kataribe
	wget https://github.com/matsuu/kataribe/releases/download/v0.4.1/kataribe-v0.4.1_linux_amd64.zip -O kataribe/kataribe.zip
	unzip -o kataribe/kataribe.zip -d kataribe
	sudo mv kataribe/kataribe /usr/local/bin/
	sudo chmod +x /usr/local/bin/kataribe
	rm -rf kataribe/
	kataribe -generate
	wget https://github.com/KLab/myprofiler/releases/download/0.2/myprofiler.linux_amd64.tar.gz
	tar xf myprofiler.linux_amd64.tar.gz
	rm myprofiler.linux_amd64.tar.gz
	sudo mv myprofiler /usr/local/bin/
	sudo chmod +x /usr/local/bin/myprofiler
	wget https://github.com/bcicen/slackcat/releases/download/v1.7.2/slackcat-1.7.2-linux-amd64 -O slackcat
	sudo mv slackcat /usr/local/bin/
	sudo chmod +x /usr/local/bin/slackcat
	slackcat --configure
	go get -u github.com/pressly/goose/cmd/goose
	go install github.com/pressly/goose/cmd/goose@latest
	cp ~/go/bin/goose /usr/local/bin/

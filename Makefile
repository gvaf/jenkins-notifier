
all: jenkins-send jenkins-notifier

jenkins-send: jenkins-send.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build jenkins-send.go

jenkins-notifier: jenkins-notifier.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build jenkins-notifier.go

run-notifier: jenkins-notifier
	mkdir -p shared/
	mkdir -p for-jenkins/
	./jenkins-notifier -socket-file ./shared/jenkins.sock -output-dir ./for-jenkins

clean: jenkins-send jenkins-notifier
	rm jenkins-send || true
	rm jenkins-notifier || true
	rm -rf for-jenkins/ || true
	rm -rf shared/ || true

guix: jenkins-send jenkins-notifier
	mkdir -p ./shared
	cp jenkins-send ./shared/
	mkdir -p for-jenkins/
	guix environment --pure --container --no-cwd --user=myuser --share=$$(pwd)/shared=/shared --ad-hoc bash python coreutils 

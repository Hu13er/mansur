#!/bin/bash


/usr/local/go/bin/go build -o ./mansur . || {
	echo "could not compile source code"
	exit 1
}

cp ./mansur /bin/mansur || {
	echo "could not cp mansur binary"
	exit 1
}

mkdir -p /etc/mansur
cp ./vars.env /etc/mansur/vars.env || {
	echo "could not cp vars.env"
	exit 1
}

cp ./mansur.service /etc/systemd/system/mansur.service || {
	echo "could not cp mansur.service"
	exit 1
}

echo "Ready to go: systemctl start mansur.service"
exit 0


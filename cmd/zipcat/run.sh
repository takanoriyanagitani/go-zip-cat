#!/bin/sh

zinput1=./sample.d/input1.zip
zinput2=./sample.d/input2.zip

zoutput=./sample.d/output.zip

geninput(){
	echo generating input zip file...

	mkdir -p ./sample.d

	echo hw1 > ./sample.d/hw1.txt
	echo hw2 > ./sample.d/hw2.txt

	echo hw3 > ./sample.d/hw3.log
	echo hw4 > ./sample.d/hw4.log

	find ./sample.d -type f -name '*.txt' |
		zip \
			-0 \
			-@ \
			-T \
			-v \
			-o \
			"${zinput1}"

	find ./sample.d -type f -name '*.log' |
		zip \
			-0 \
			-@ \
			-T \
			-v \
			-o \
			"${zinput2}"

}

test -f "${zinput1}" || geninput
test -f "${zinput2}" || geninput

echo concatenating zip files...
ls "${zinput1}" "${zinput2}" |
	./zipcat |
	dd \
		if=/dev/stdin \
		of="${zoutput}" \
		bs=1048576 \
		status=none \
		conv=fsync

ls -lSh \
	"${zinput1}" \
	"${zinput2}" \
	"${zoutput}"

unzip -lv "${zinput1}"
unzip -lv "${zinput2}"
unzip -lv "${zoutput}"

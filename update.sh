#!/bin/bash
## DeGOps: 0.0.4
# Updates DeGOps files on current dir.
set -o errexit
set -o nounset

## Validate action

getmake() {
    echo "  Load remote Makefile"
	curl -o Makefile https://raw.githubusercontent.com/jimmy-go/degops/develop/Makefile
}

getdocker() {
    echo "  Load remote Dockerfile"
	curl -o Dockerfile https://raw.githubusercontent.com/jimmy-go/degops/develop/Dockerfile
}

gettravis() {
    echo "  Load remote .travis.yml"
	curl -o .travis.yml https://raw.githubusercontent.com/jimmy-go/degops/develop/.travis.yml
}

getupdate() {
    echo "  Load remote update script"
	curl -o update.sh https://raw.githubusercontent.com/jimmy-go/degops/develop/update.sh
}

getscript() {
	FILE=$1

	# Check for scripts overwrites.
	OW="${FILE}_overwrite.sh"
	if [ -f $OW ]; then
		echo "File overwrite found, skip: $OW"
		return
	fi

	# Load remote file.
    echo "  Load remote script: $FILE"
	curl -o scripts/${FILE}.sh https://raw.githubusercontent.com/jimmy-go/degops/develop/scripts/${FILE}.sh
    chmod +x scripts/${FILE}.sh
}

copyscripts() {
	mkdir -p scripts
	getscript "install"
	getscript "test"
	getscript "run"
	getscript "cover"
	getscript "clean"
}

case "$1" in
	all)
		getupdate
		getmake
		copyscripts
		getdocker
		gettravis
		;;
	update)
		getupdate
		;;
	makefile)
		getmake
		;;
	scripts)
		copyscripts
		;;
	container)
		getdocker
		;;
	travis)
		gettravis
		;;
	*)
		echo $"Usage: $0 {all|scripts|makefile|container|update|travis}"
		exit 1
esac

## Download.

## Copy overwrites.

# curl -o scripts/test.sh https://raw.githubusercontent.com/jimmy-go/degops/develop/scripts/test.sh

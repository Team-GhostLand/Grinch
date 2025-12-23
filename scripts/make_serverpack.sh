#!/bin/bash

if [ -z "$1" ]; then
    echo "Pack name missing!";
    exit 1;
fi


mkdir ".temp" || exit
cd "./.temp/" || exit
unzip "../$1" || exit
cd "./overrides/" || exit

RMINDEX="../server-overrides/REMOVALS.txt"
if [ -f "$RMINDEX" ]; then
	IFS=$'\n'
	# shellcheck disable=SC2046
	rm -rv $(cat "$RMINDEX") || exit
	
	rm "$RMINDEX" || exit
fi
cd .. || exit

zip -r "serverpack-$1" ./* || exit
mv "./serverpack-$1" "./../serverpack-$1" || exit
cd ..
rm -r "./.temp/" || exit
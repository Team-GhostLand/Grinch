#!/bin/bash

if [ -z "$1" ]; then
    echo "Pack name missing!";
    exit 1;
fi


mkdir ".temp" || exit
cd "./.temp/" || exit
unzip "../$1" || exit
cd "./overrides/" || exit

IFS=$'\n'
# shellcheck disable=SC2046
rm -rv $(cat "../server-overrides/REMOVALS.txt") || exit

rm "../server-overrides/REMOVALS.txt" || exit
cd .. || exit

zip -r "serverpack-$1" ./* || exit
mv "./serverpack-$1" "./../serverpack-$1" || exit
cd ..
rm -r "./.temp/" || exit
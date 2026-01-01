#!/bin/sh

# THIS IS SUPPOSED TO ONLY BE RUN FROM A DOCKER IMAGE!

echo "
-----------STARTED ON: $(date)-----------";

if [ -z "$REPO" ]; then
	if [ -z "$REPO_FILE" ]; then
    	echo "ERROR: You must specify a \$REPO envar!";
    	exit 1
	else
		REPO="$(cat "$REPO_FILE")"
	fi
fi
if [ -z "$PACK" ]; then
	if [ -z "$PACK_FILE" ]; then	
		echo "ERROR: You must specify a \$PACK envar!";
		exit 1
	else
		PACK="$(cat "$PACK_FILE")"
	fi
fi
echo "Using repo: $REPO";
echo "Using pack: $PACK";

if [ -e "cache" ]; then
    echo "...which is already cached - will do a simple git pull to check for updates instead of a full git clone.";
    cd "cache" || exit
    git pull || exit
else
    echo "...which has to be cloned.";
    git clone "$REPO" "cache" || exit
    echo "[STARTING FROM SCRATCH]" > "last-version.txt"
    cd "cache" || exit
fi

VER="$(/app/grinch query version "$PACK")"

if [ "$VER" = "$(cat ../last-version.txt)" ]; then
    echo "You seem to be using an up-to date modpack. Waiting for 15s until the next cycle.";
    sleep 15
    cd "/workdir" || exit
    exec "/app/ci.sh";
    echo "ERROR: If you're seeing this, the next cycle couldn't be started!";
    exit 1;
fi

echo "Building modpack version $VER (because it's different than $(cat ../last-version.txt))";
MRP=".mrpack"
MAIN="$(/app/grinch query name "$PACK")$MRP"
S="$(/app/grinch query name_slim "$PACK")$MRP"
T="$(/app/grinch query name_tweakable "$PACK")$MRP"
/app/grinch e -T "$MAIN" "$PACK" || exit
/app/grinch e -sT "$S" "$PACK" || exit
/app/grinch e -tT "$T" "$PACK" || exit

echo "Exporting just-built assets...";
mv ./*.mrpack "/exports" || exit
echo "$VER" > "../last-version.txt"

UPDATE_TARGET="latest_server.mrpack"
echo "Updating $UPDATE_TARGET";
cd "/exports" || exit
rm "$UPDATE_TARGET";
if [ $? -ne 0 ]; then
    echo "WARN: Failed to delete the older $UPDATE_TARGET. Not treating it as an error, as it simply must've already not existed. (And if it does exist - the next command will either throw an error and properly exit or simply overwrite it.)";
fi
cp "$MAIN" "$UPDATE_TARGET" || exit

echo "DONE! Waiting for 2min until the next cycle.";
sleep 120
cd "/workdir" || exit
exec "/app/ci.sh";
echo "ERROR: If you're seeing this, the next cycle couldn't be started!";
exit 1;
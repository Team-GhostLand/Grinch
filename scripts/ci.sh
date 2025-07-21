#!/bin/sh

# THIS IS SUPPOSED TO ONLY BE RUN FROM A DOCKER IMAGE!

echo "-----------STARTED ON: $(date)-----------";

if [ -z "$REPO" ]; then
    echo "ERROR: You must specify a \$REPO envar!";
    exit 1
fi
if [ -z "$NAME" ]; then
    echo "ERROR: You must specify a \$NAME envar!";
    exit 1
fi
echo "Using repo: $REPO";

if [ -e "cache" ]; then
    echo "...which is already cached - will do a simple git pull to check for updates instead of a full git clone.";
    cd "cache" || exit
    git pull || exit
else
    echo "...which has to be cloned.";
    git clone "$REPO" "cache" || exit
    echo "[STARTING FROM SCRATCH]" >> "last-version.txt"
    cd "cache" || exit
fi

VER="$(/app/grinch vq)"

if [ "$VER" = "$(cat ../last-version.txt)" ]; then
    echo "You seem to be using an up-to date modpack. Waiting for 15s until the next cycle.";
    sleep 15
    exec "/app/ci.sh";
    echo "ERROR: If you're seeing this, the next cycle couldn't be started!";
    exit 1;
fi

echo "Building modpack version $VER";
MRP=".mrpack"
Q="quick$MRP"
S="slim$MRP"
T="tweakable$MRP"
/app/grinch e -qT "$Q" || exit
/app/grinch e -sT "$S" || exit
/app/grinch e -tT "$T" || exit
/app/grinch-serverpack "$Q" || exit

EXPORTNAME="$NAME $(/app/grinch vq)";
echo "Exporting just-built assets as $EXPORTNAME";
SERVERPACK="$EXPORTNAME - Server Edition$MRP"
mv "$Q" "$EXPORTNAME$MRP" || exit
mv "$S" "$EXPORTNAME - Slim Edition$MRP" || exit
mv "$T" "$EXPORTNAME - Tweakable Edition$MRP" || exit
mv "serverpack-$Q" "$SERVERPACK" || exit
mv ./*.mrpack "/exports" || exit
echo "$VER" >> "last-version.txt"

UPDATE_TARGET="latest_server.mrpack"
echo "Updating $UPDATE_TARGET";
cd "/exports" || exit
rm "$UPDATE_TARGET";
if [ $? -ne 0 ]; then
    echo "WARN: Failed to delete the older $UPDATE_TARGET. Not treating it as an error, as it simply must've already not existed. (And if it does exist - the next command will either throw an error and properly exit or simply overwrite it.)";
fi
cp "$SERVERPACK" "$UPDATE_TARGET" || exit

echo "DONE! Waiting for 2min until the next cycle.";
sleep 120
exec "/app/ci.sh";
echo "ERROR: If you're seeing this, the next cycle couldn't be started!";
exit 1;
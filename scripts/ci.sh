#!/bin/sh

# THIS IS SUPPOSED TO ONLY BE RUN FROM A DOCKER IMAGE!

echo "STARTED ON: $(date)";
echo "This is ATM just a simple test whether Docker works!!!";
if [ -z "$REPO" ]; then
    echo "ERROR: You must specify a \$REPO envar!";
    exit 1
fi
echo "Uses repo: $REPO";
echo "Grinch's help outupt:";
./grinch
echo "...and Git's:";
git help
echo "......and make_serverpacks's (should complain about a missing name):";
./make_serverpack.sh
echo "Test completed. This command will soon exit and (depending on your setup) might be restarted by Docker.";
sleep 5;
exit 0;
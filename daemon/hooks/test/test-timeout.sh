#!/bin/sh -eu

echo "TEST $ZREPL_HOOKTYPE $ZREPL_FS@$ZREPL_SNAPNAME ZREPL_TIMEOUT=$ZREPL_TIMEOUT"

exec sleep $(($ZREPL_TIMEOUT + 1))
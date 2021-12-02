#!/bin/sh

# https://github.com/traviscross/mtr/issues/415

# $1 = host ip
# $2 = port

mtr -T $1 -P $2 --report-cycles 10 -n --csv >> "$PEER_OUT_DIR/$1:$2.csv"

# mtr -T $1 -P $2 --report-cycles 10 -n --csv # TODO: report-wide-csv to capture secondary servers

#!/bin/bash
echo "start"
./timingGetData &
./luckyGd >> luckyGd.log &
./luckyJx >> luckyJx.log
echo "end"
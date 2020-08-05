#!/bin/bash
echo "start"
./timingGetData &
./lucky &
./botmsg
echo "end"
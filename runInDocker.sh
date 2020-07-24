#!/bin/bash
echo "start"
./timingGetData &
./lucky &
./botMsg
echo "end"
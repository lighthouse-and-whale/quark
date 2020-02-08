#!/usr/bin/env bash
kill $(cat pid)
nohup ./server >> nohup.log & printf $! > pid
sleep 0.1s

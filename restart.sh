#!/bin/sh
ps aux | grep logseq-auto-push | grep -v grep | awk '{print $2}' | xargs kill

nohup ./logseq-auto-push &

#!/bin/bash
curl \
    -s \
    http://localhost:3000/ | jq .
    #http://localhost:3000/proxy/?q=https://www.skuid.com | jq .

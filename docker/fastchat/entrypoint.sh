#!/bin/bash

nohup python3 -m fastchat.serve.controller --host 0.0.0.0 --port 21001 > /dev/null 2>&1 &

# Wait for the controller to start
sleep 5

# Start the openapi server
python3 -m fastchat.serve.openai_api_server --host 0.0.0.0 --port 8000 \
--controller-address http://$(hostname -I | awk '{print $1}'):21001
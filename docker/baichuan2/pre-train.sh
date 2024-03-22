#!/bin/bash
hostfile=''
deepspeed -hostfile=$hostfile \
--force_multi \
train.py \
--deepspeed \
--deepspeed_config deepspeed.json
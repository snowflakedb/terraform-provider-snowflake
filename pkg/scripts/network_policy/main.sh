#!/bin/bash

# Filter out only IPv4 CIDR lines from test.txt and overwrite test.txt
curl https://api.github.com/meta | jq -r '.actions[]' |grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+(/[0-9]+)?$' > github_ipv4_subnets.txt

terraform plan

#!/bin/bash

set -e

# Pre-install pacapt
# ADD https://github.com/icy/pacapt/raw/ng/pacapt /usr/local/bin/pacapt
curl -sSL https://github.com/icy/pacapt/raw/ng/pacapt \
  -o /usr/local/bin/pacapt && \
  chmod 0755 /usr/local/bin/pacapt

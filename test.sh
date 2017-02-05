#!/bin/bash

cl() { ls -la | wc -l; }

cd
# pwd
c1=$(cl)
dot -R $DOT
c2=$(cl)
dot -R $DOT remove
c3=$(cl);

# echo "$c1 -install-> $c2 -remove-> $c3"

[[ "$c1" -eq "$c3" ]] \
  && [[ "$c2" -gt "$c1" ]] \
  && [[ "$c2" -gt "$c3" ]]

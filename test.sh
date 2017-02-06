#!/bin/bash

cl() { ls -la | wc -l; }

cd
# pwd
c1=$(cl) \
&& t1=$(tail -n 1 ~/.bashrc) \
&& dot -R $DOT --non-interactive \
&& c2=$(cl) \
&& t2=$(tail -n 1 ~/.bashrc) \
&& dot -R $DOT remove --non-interactive \
&& c3=$(cl) \
&& t3=$(tail -n 1 ~/.bashrc)

# echo "$c1 -install-> $c2 -remove-> $c3"

[[ "$c1" -eq "$c3" ]] && [[ "$c2" -gt "$c1" ]] && [[ "$c2" -gt "$c3" ]] \
&& [[ "$t1" == "$t3" ]] && [[ "$t1" != "$t2" ]] && [[ "$t2" != "$t3" ]]

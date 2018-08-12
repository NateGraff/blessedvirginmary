#!/usr/bin/env bash

# Need Bash >= 4.x for associative arrays
declare -A TEST_LIST=(
        [basic/ret_zero.ll]="0"
        [basic/array.ll]="15"
        [basic/ptr.ll]="5"
        [basic/ptrtoarray.ll]="0"
# math tests
        [math/oneplusone.ll]="2"
        [math/twotimestwo.ll]="4"
        [math/fourdivtwo.ll]="2"
        [math/threeminusone.ll]="2"
        [math/math_stuff.ll]="0"
        [math/comparison.ll]="5"

# branching tests
        [branch/branch.ll]="2"
        [branch/falsebranch.ll]="3"
        [branch/nestedbranch.ll]="2"

# function tests
        [functions/noargs.ll]="1"
        [functions/playnice.ll]="1"
#        [addone.ll]="2"
)

echo "-------------------------"
echo "Testing blessedvirginmary"
echo "-------------------------"

parallel --link ./runtest.bash ::: ${!TEST_LIST[@]} ::: ${TEST_LIST[@]}

echo "-------------------------"

exit ${FAIL_NUM}


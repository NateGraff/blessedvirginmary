#!/usr/bin/env bash

COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_NC='\033[0m'

# Need Bash >= 4.x for associative arrays
declare -A TEST_LIST=(
        [ret_zero]="0"
        [array]="15"
        [math_stuff]="0"
        [ptr]="5"
        [ptrtoarray]="0"
)

PASS_NUM=0
FAIL_NUM=0

echo "-------------------------"
echo "Testing blessedvirginmary"
echo "-------------------------"

for TEST_NAME in "${!TEST_LIST[@]}"; do
        TEST_EXPECTED=${TEST_LIST[$TEST_NAME]}

        # make the LLVM IR
        make -s -C reference ${TEST_NAME}.ll

        # execute the transpiled Bash
        ./blessedvirginmary reference/${TEST_NAME}.ll | bash

        TEST_RETURN=$?
        if [[ $TEST_EXPECTED -eq $TEST_RETURN ]]; then
                echo -e " ${COLOR_GREEN}PASS ${COLOR_NC}${TEST_NAME}"
                PASS_NUM=$(expr ${PASS_NUM} + 1)
        else
                echo -e " ${COLOR_RED}FAIL ${COLOR_NC}${TEST_NAME}"
                echo "    Expected $TEST_EXPECTED but got $TEST_RETURN"
                FAIL_NUM=$(expr ${FAIL_NUM} + 1)
        fi
done

echo "-------------------------"
echo -e "Pass count: ${COLOR_GREEN}${PASS_NUM}${COLOR_NC} Fail count: ${COLOR_RED}${FAIL_NUM}${COLOR_NC}"

exit ${FAIL_NUM}


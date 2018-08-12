#!/usr/bin/env bash

COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_NC='\033[0m'

testoutput() {
        TEST_RETURN=$1
        TEST_EXPECTED=$2
        if [[ $TEST_EXPECTED -eq $TEST_RETURN ]]; then
                echo -e " ${COLOR_GREEN}PASS ${COLOR_NC}${TEST_NAME}"
                return 0
        else
                echo -e " ${COLOR_RED}FAIL ${COLOR_NC}${TEST_NAME}"
                echo "    Expected $TEST_EXPECTED but got $TEST_RETURN"
                return 1
        fi
}

runtest() {
        TEST_NAME=$1
        TEST_EXPECTED=$2
        make -s -C reference ${TEST_NAME}
        ./blessedvirginmary reference/${TEST_NAME} | bash
        testoutput $? ${TEST_EXPECTED}
        return $?
}

runtest $1 $2
exit $?

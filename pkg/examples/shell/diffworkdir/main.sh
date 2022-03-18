#!/bin/bash
# Linked to https://app.airplane.dev/t/new_shell_task [do not edit this line]
# Params are in environment variables as PARAM_{SLUG}, e.g. PARAM_USER_ID

# test.sh is created in the custom Dockerfile in directory /myScripts.
# Testing this ensures that we execute the main shell script within the same
# working directory as the custom Dockerfile.
./test.sh

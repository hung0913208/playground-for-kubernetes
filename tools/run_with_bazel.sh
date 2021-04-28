#!/bin/bash

######################################################################
# @author      : hung0913208 (hung0913208@gmail.com)
# @file        : run_with_bazel
# @created     : Tuesday Apr 27, 2021 20:59:28 +07
#
# @description : Run bazel command to do several stub
######################################################################

CMD=$1

shift

if [ "$CMD" = "test" ]; then
	PARAMS="test --test_tag_filters=selftest $@"
elif [ "$CMD" = "build" ]; then
	PARAMS="build $@"
elif [ "$CMD" = "gazele" ]; then
	PARAMS="run //:gazelle"
elif [ "$CMD" = "update-repos" ]; then
	PARAMS="run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies"
fi

if which docker &> /dev/null && [[ ${#PARAMS} -gt 0 ]]; then
	docker run 				\
		-e USER="root" 			\
		-u="$(id -u)" 			\
		-v $(pwd):$(pwd) 		\
		-v $(pwd):$(pwd) 		\
		-w $(pwd) 			\
		l.gcr.io/google/bazel:latest 	\
		--output_user_root=$(pwd) 	\
		$PARAMS
fi

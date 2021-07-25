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
	PARAMS="test --test_tag_filters=selftest --test_output=errors $@"
elif [ "$CMD" = "build" ]; then
	PARAMS="build $@"
elif [ "$CMD" = "gazelle" ]; then
	PARAMS="run //:gazelle"
elif [ "$CMD" = "update-repos" ]; then
	PARAMS="run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies"
fi

if [ ! -f ./go.sum ]; then
	touch ./go.sum

	if which go &> /dev/null; then
		go list -m -json all
	elif which docker &> /dev/null; then
		docker run 								\
			-e USER="$(id -u)" 			\
			-u="$(id -u)" 					\
			-v $(pwd):$(pwd) 				\
			-v $(pwd):$(pwd) 				\
			-w $(pwd) 							\
			golang go list -m -json all
	else
		exit -1
	fi
fi

if which bazel &> /dev/null; then
	bazel $PARAMS
elif which docker &> /dev/null && [[ ${#PARAMS} -gt 0 ]]; then
	IMAGE="l.gcr.io/google/bazel:2.2.0"

	docker run 										\
		-e USER="$(id -u)" 					\
		-u="$(id -u)" 							\
		-v $(pwd):$(pwd) 						\
		-v $(pwd):$(pwd) 						\
		-w $(pwd) 									\
		$IMAGE 											\
		--output_user_root=$(pwd) 	\
		$PARAMS
fi

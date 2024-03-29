#!/bin/sh
if [[ -z "$GIT_REPO" ]]; then
    echo "error: git repo not defined"
    exit 1
fi

if [[ -z "$VERSION" ]]; then
    echo "error: missing VERSION"
    exit 1
fi

if [[ -z "$GITHUB_TOKEN" ]]; then
    echo "error: GITHUB_TOKEN token not defined"
    exit 1
fi

if [[ -z "$PRERELEASE" ]]; then
    PRERELEASE=false
fi


if [[ "$DRY_RUN" == "1" ]]; then
    echo "gh release create --repo $GIT_REPO -t $VERSION -n $VERSION --prerelease=$PRERELEASE $VERSION"
    exit 0
fi

gh release create --repo $GIT_REPO -t $VERSION -n $VERSION --prerelease=$PRERELEASE $VERSION

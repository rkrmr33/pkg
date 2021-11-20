#!/bin/sh

echo "checking worktree..."
res=$(git status -s)
if [[ -z "$res" ]]; then
    echo clean worktree
    exit 0
fi


GIT_PAGER=cat git diff --minimal
echo
echo ERROR: worktree is not clean! make sure you run \"make gen\" and commit the changes before.
echo
exit 1
#!/usr/bin/env bash

echo "Running version bump script"

function forward ()
{
  VERSION=$(cat ./version)
  NEW_VERSION=$(echo $VERSION | awk -F. '/[0-9]+\./{$NF++;print}' OFS=.)
  echo $NEW_VERSION > ./version
}


function revert () {
  VERSION=$(cat ./version)
  NEW_VERSION=$(echo $VERSION | awk -F. '/[0-9]+\./{$NF--;print}' OFS=.)
  echo $NEW_VERSION > ./version
}

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)


if [[ "$CURRENT_BRANCH" != "main" ]]; then
  echo "You are not on the main branch, it will not properly release if you're not on the main branch"
  exit 1
fi

git fetch
git pull origin main

if [[ $? != 0 ]]; then
  echo "Could pull origin main, make sure you don't have lingering local changes or diverging history"
  exit 1
fi

forward
git tag $(cat ./version)

if [[ $? != 0 ]]; then
  echo "failed to tag commit, reverting version"
  revert
  exit 1
fi

git push origin $(cat ./version)

if [[ $? != 0 ]]; then
  echo "failed to push tag, reverting version"
  revert
  exit 1
fi


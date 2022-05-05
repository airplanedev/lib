# Contributing to the Airplane lib

## Tagging

To tag a new version of the Airplane lib, run the following command with the new tag you want to create:

```sh
export AIRPLANE_LIB_TAG=v0.0.1 && \
  git tag ${AIRPLANE_LIB_TAG} && \
  git push origin ${AIRPLANE_LIB_TAG}
```

## Base images

If you make changes to the base images -- `pkg/build/versions.json` -- then you'll need to push these into the public cache:

```
airplane dev scripts/cache/main.ts -- --gcp_project=airplane-stage
airplane dev scripts/cache/main.ts -- --gcp_project=airplane-prod
```

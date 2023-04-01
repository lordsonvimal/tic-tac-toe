#!/usr/bin/env bash
pnpm -C client/ dev & gin --appPort 8080 --path server/

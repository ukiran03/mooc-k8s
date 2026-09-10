#!/usr/bin/env bash

# NOTE: MUST Run this script from "scripts/" directory itself
kubectl apply -f ../base/app/secrets.yaml -n production

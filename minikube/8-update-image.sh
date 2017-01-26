#!/bin/bash

kubectl set image deployment/lunchtalk lunchtalk=micahhausler/k8s-lunchtalk:latest

kubectl get deployment lunchtalk -o yaml

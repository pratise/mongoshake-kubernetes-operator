#!/bin/bash
kubectl delete -f ./manager/manager.yaml -n operator
kubectl delete -f ./rbac/role_binding.yaml -n operator
kubectl delete -f ./rbac/service_account.yaml -n operator
kubectl delete -f ./rbac/role.yaml -n operator
kubectl delete -f ./crd/bases/pratise.github.com_mongoshakes.yaml




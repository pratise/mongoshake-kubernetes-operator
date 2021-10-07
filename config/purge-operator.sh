#!/bin/bash
kubectl delete -f ./manager/manager.yaml -n operator
kubectl delete -f ./manager/manager_config.yaml -n operator
kubectl delete -f ./rbac/leader_election_role_binding.yaml -n operator
kubectl delete -f ./rbac/role_binding.yaml -n operator
kubectl delete -f ./rbac/service_account.yaml -n operator
kubectl delete -f ./rbac/leader_election_role.yaml -n operator
kubectl delete -f ./rbac/mongoshake_viewer_role.yaml -n operator
kubectl delete -f ./rbac/mongoshake_editor_role.yaml -n operator
kubectl delete -f ./rbac/role.yaml -n operator
kubectl delete -f ./crd/bases/pratise.github.com_mongoshakes.yaml

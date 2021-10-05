#!/bin/bash
kubectl apply -f ./crd/bases/pratise.github.com_mongoshakes.yaml
kubectl apply -f ./rbac/role.yaml -n operator
kubectl apply -f ./rbac/leader_election_role.yaml -n operator
kubectl apply -f ./rbac/mongoshake_editor_role.yaml -n operator
kubectl apply -f ./rbac/mongoshake_viewer_role.yaml -n operator
kubectl apply -f ./rbac/service_account.yaml -n operator
kubectl apply -f ./rbac/role_binding.yaml -n operator
kubectl apply -f ./rbac/leader_election_role_binding.yaml -n operator
kubectl apply -f ./manager/manager_config.yaml -n operator
kubectl apply -f ./manager/manager.yaml -n operator

#!/bin/sh
{ 
    # コネクション1 の開始
    (curl http://$(kubectl get svc signal -ojsonpath="{.status.loadBalancer.ingress[0].ip}")/client1 &)
    # 5 秒間の待機
    sleep 5
    # Pod の削除
    (kubectl delete po $(kubectl get po -ojsonpath="{.items[0].metadata.name}") &)
    # 1 秒間の待機
    sleep 1
    # コネクション2 の開始
    (curl http://$(kubectl get svc signal -ojsonpath="{.status.loadBalancer.ingress[0].ip}")/client2 &)
}

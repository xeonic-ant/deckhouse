---
title: "Модуль user-authz"
---

Модуль отвечает за генерацию RBAC для пользователей и реализует простейший режим multi-tenancy с разграничением доступа по namespace.

Реализует role-based-подсистему сквозной авторизации, расширяя функционал стандартного механизма RBAC.

Вся настройка прав доступа происходит с помощью [Custom Resources](cr.html).

## Возможности модуля

- Управление доступом пользователей и групп на базе механизма RBAC Kubernetes
- Управление доступом к инструментам масштабирования (параметр `allowScale` Custom Resource [`ClusterAuthorizationRule`](cr.html#clusterauthorizationrule))
- Управление доступом к форвардингу портов (параметр `portForwarding` Custom Resource [`ClusterAuthorizationRule`](cr.html#clusterauthorizationrule))
- Управление списком разрешенных namespace в формате регулярных выражений (параметр `limitNamespaces` Custom Resource [`ClusterAuthorizationRule`](cr.html#clusterauthorizationrule))
- Управление доступом к системным namespace (параметр `allowAccessToSystemNamespaces` Custom Resource [`ClusterAuthorizationRule`](cr.html#clusterauthorizationrule)), таким как `kube-system` и пр;

## Ролевая модель

В модуле кроме прямого использования RBAC можно использовать удобный набор высокоуровневых ролей:
- `User` — позволяет получать информацию обо всех объектах (включая доступ к журналам Pod'ов), но не позволяет заходить в контейнеры, читать секреты и выполнять port-forward;
- `PrivilegedUser` — то же самое, что и `User`, но позволяет заходить в контейнеры, читать секреты, а также позволяет удалять Pod'ы (что обеспечивает возможность перезагрузки);
- `Editor` — то же самое, что и `PrivilegedUser`, но предоставляет возможность создавать, изменять и удалять все объекты, которые обычно нужны для прикладных задач.
- `Admin` — то же самое, что и `Editor`, но позволяет удалять служебные объекты (производные ресурсы, например, `ReplicaSet`, `certmanager.k8s.io/challenges` и `certmanager.k8s.io/orders`);
- `ClusterEditor` — то же самое, что и `Editor`, но позволяет управлять ограниченным набором `cluster-wide` объектов, которые могут понадобиться для прикладных задач (`ClusterXXXMetric`, `KeepalivedInstance`, `DaemonSet` и т.д). Роль для работы оператора кластера.
- `ClusterAdmin` — то же самое, что и `ClusterEditor` + `Admin`, но позволяет управлять служебными cluster-wide объектами (производные ресурсы, например, `MachineSets`, `Machines`, `OpenstackInstanceClasses`..., а так же `ClusterAuthorizationRule`, `ClusterRoleBindings` и `ClusterRole`). Роль для работы администратора кластера. **Важно**, что так как `ClusterAdmin` уполномочен редактировать `ClusterRoleBindings`, он может **сам себе расширить полномочия**.
- `SuperAdmin` — разрешены любые действия с любыми объектами, при этом ограничения [`limitNamespaces`](#возможности-модуля) продолжат работать.

## Особенности реализации

> **Важно!** Режим multi-tenancy (авторизация по namespace) в данный момент реализован по временной схеме и **не гарантирует безопасность!**

Если webhook, который реализовывает систему авторизации по какой-то причине будет недоступен, то в это время опции `allowAccessToSystemNamespaces` и `limitNamespaces` в CR перестанут применяться и пользователи будут иметь доступ во все namespace. После восстановления доступности webhook'а опции продолжат работать.

## Список доступа для каждой роли модуля по умолчанию

Сокращения для `verbs`:
<!-- start user-authz roles placeholder -->
* read - `get`, `list`, `watch`
* read-write - `get`, `list`, `watch`, `create`, `delete`, `deletecollection`, `patch`, `update`
* write - `create`, `delete`, `deletecollection`, `patch`, `update`

{{site.data.i18n.common.role[page.lang] | capitalize }} `User`:

```text
read:
    - apiextensions.k8s.io/customresourcedefinitions
    - apps/daemonsets
    - apps/deployments
    - apps/replicasets
    - apps/statefulsets
    - autoscaling.k8s.io/verticalpodautoscalers
    - autoscaling/horizontalpodautoscalers
    - batch/cronjobs
    - batch/jobs
    - configmaps
    - discovery.k8s.io/endpointslices
    - endpoints
    - events
    - events.k8s.io/events
    - extensions/daemonsets
    - extensions/deployments
    - extensions/ingresses
    - extensions/replicasets
    - extensions/replicationcontrollers
    - limitranges
    - metrics.k8s.io/nodes
    - metrics.k8s.io/pods
    - namespaces
    - networking.k8s.io/ingresses
    - networking.k8s.io/networkpolicies
    - nodes
    - persistentvolumeclaims
    - persistentvolumes
    - pods
    - pods/log
    - policy/poddisruptionbudgets
    - rbac.authorization.k8s.io/rolebindings
    - rbac.authorization.k8s.io/roles
    - replicationcontrollers
    - resourcequotas
    - serviceaccounts
    - services
    - storage.k8s.io/storageclasses
```

{{site.data.i18n.common.role[page.lang] | capitalize }} `PrivilegedUser` ({{site.data.i18n.common.includes_rules_from[page.lang]}} `User`):

```text
create:
    - pods/eviction
create,get:
    - pods/attach
    - pods/exec
delete,deletecollection:
    - pods
read:
    - secrets
```

{{site.data.i18n.common.role[page.lang] | capitalize }} `Editor` ({{site.data.i18n.common.includes_rules_from[page.lang]}} `User`, `PrivilegedUser`):

```text
read-write:
    - apps/deployments
    - apps/statefulsets
    - autoscaling.k8s.io/verticalpodautoscalers
    - autoscaling/horizontalpodautoscalers
    - batch/cronjobs
    - batch/jobs
    - configmaps
    - discovery.k8s.io/endpointslices
    - endpoints
    - extensions/deployments
    - extensions/ingresses
    - networking.k8s.io/ingresses
    - persistentvolumeclaims
    - policy/poddisruptionbudgets
    - serviceaccounts
    - services
write:
    - secrets
```

{{site.data.i18n.common.role[page.lang] | capitalize }} `Admin` ({{site.data.i18n.common.includes_rules_from[page.lang]}} `User`, `PrivilegedUser`, `Editor`):

```text
create,patch,update:
    - pods
delete,deletecollection:
    - apps/replicasets
    - extensions/replicasets
```

{{site.data.i18n.common.role[page.lang] | capitalize }} `ClusterEditor` ({{site.data.i18n.common.includes_rules_from[page.lang]}} `User`, `PrivilegedUser`, `Editor`):

```text
read:
    - rbac.authorization.k8s.io/clusterrolebindings
    - rbac.authorization.k8s.io/clusterroles
write:
    - apiextensions.k8s.io/customresourcedefinitions
    - apps/daemonsets
    - extensions/daemonsets
    - storage.k8s.io/storageclasses
```

{{site.data.i18n.common.role[page.lang] | capitalize }} `ClusterAdmin` ({{site.data.i18n.common.includes_rules_from[page.lang]}} `User`, `PrivilegedUser`, `Editor`, `Admin`, `ClusterEditor`):

```text
read-write:
    - deckhouse.io/clusterauthorizationrules
write:
    - limitranges
    - namespaces
    - networking.k8s.io/networkpolicies
    - rbac.authorization.k8s.io/clusterrolebindings
    - rbac.authorization.k8s.io/clusterroles
    - rbac.authorization.k8s.io/rolebindings
    - rbac.authorization.k8s.io/roles
    - resourcequotas
```
<!-- end user-authz roles placeholder -->

Вы можете получить дополнительный список правил доступа для роли модуля из кластера ([существующие пользовательские правила](usage.html#customizing-rights-of-high-level-roles) и не стандартные правила из других модулей Deckhouse):

```bash
D8_ROLE_NAME=Editor
kubectl get clusterrole -A -o jsonpath="{range .items[?(@.metadata.annotations.user-authz\.deckhouse\.io/access-level=='$D8_ROLE_NAME')]}{.rules}{'\n'}{end}" | jq -s add
```

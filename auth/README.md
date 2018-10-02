# Access Control List ( ACL )

#### Description:
ACL can be used to authorize, restrict or reject certain access, action or usage to resources. 
 
#### Use Case:
ACL is highly useful for putting restrictions on user capabilities such as assigning roles and dividing them into groups with certain permissions. ACL is divided into four categories, **entity**, **ability**, **action** and **resource**, in this particular order.  Entity refers to a group or a user that is initiating some action on a resource. Ability is a binary **Can** or **Cannot** identifier used to indicate whether a particular user has the sufficient permission to proceed. An action is an identifier that indicates what kind of operation is going to take place. The resource is the final block in the chain which evaluates the ability of the user and kind of requested action on itself in order to accept or refuse accesses and actions. Subsequently, these properties can be inherited by subgroups. 

For example, an application with three user groups can be illustrated as following.

| Guest  | User | Admin |
| --- | --- | --- |


| Entity | Ability | Action | Resource
| --- | --- | --- | --- |
| User | *can* | subscribe to | $channel |
| User | *cannot* | publish to | $channel |

**Access Types** 

There are four major access types which are a combination of two main types, **Inclusive** and **Exclusive**. Inclusive allows accesses to anything that is explicitly defined and rejects anything that is undefined. Exclusive allows accesses to anything that is explicitly undefined and rejects anything that is defined ( opposite of Inclusive ). With this two definitions, following schema can be derived.

| Type | Description |
| --- | --- |
| Inclusive, Exclusive | can use anything that is DEFINED, rejects anything that is undefined. |
| Exclusive, Inclusive | can use anything that is UNDEFINED, rejects anything that is DEFINED. |
| Inclusive, Inclusive | can use anything that is only DEFINED. |
| Exclusive, Exclusive | can use anything that is only UNDEFINED. |

Differences between this "fixed" version of the gitea spec and the original one:
- at the moment, we consider that elements given in #/responses do not constitue valid schema object
- the "name": "string" constructs causes a name conflict with the String() method when urlbuilder is generated (may be disabled)

649c649,652
<             "$ref": "#/responses/empty"
---
>             "description": "user is a member",
>             "schema": {
>               "$ref": "#/responses/empty"
>             }
652c655,658
<             "$ref": "#/responses/empty"
---
>             "description": "user is not a member",
>             "schema": {
>               "$ref": "#/responses/empty"
>             }
683c689,692
<             "$ref": "#/responses/empty"
---
>             "description": "member removed",
>             "schema": {
>               "$ref": "#/responses/empty"
>             }
739c748,751
<              "$ref": "#/responses/empty"
---
>             "description": "user is a public member",
>             "schema": {
>               "$ref": "#/responses/empty"
>             }
741a754,755
>             "description": "user is not a public member",
>             "schema": {
742a757
>             }
772a788,789
>             "description": "membership publicized",
>             "schema": {
773a791
>             }
1832c1850
<             "name": "myString",
---
>             "name": "string",
2145c2163
<             "name": "myString",
---
>             "name": "string",
3301a3320,3321
>             "description": "pull request has been merged",
>             "schema": {
3302a3323
>             }
3304a3326,3327
>             "description": "pull request has not been merged",
>             "schema": {
3305a3329
>             }
4206a4231,4232
>             "description": "team deleted",
>             "schema": {
4207a4234
>             }

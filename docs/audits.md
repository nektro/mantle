# Audits

Full database dump can be found at `http://MATLE_PATH/api/admin/audits.csv`.

Using audits, administrators are able to get a full view of the entire history of different actions taken place on their Mantle server. All audits are held in a single append-only table. Descriptions and Keys on how to read this table and its output are described below.

## Schema
| Field | Type | Description |
|-------|------|-------------|
| ID | `uint` | Canonically ordered and incremented with new actions |
| ULID | `ULID` | https://github.com/ulid/spec |
| Created On | `string` | RFC 3339 date-time |
| Action | `uint` | enum described below |
| Agent | `ULID` | User who performed the action |
| Affected | `ULID` | Resource affected by action |
| Key | `string` | Specific to each action type |
| Value | `string` | Specific to each action type |

## Actions
| Number | Name |
|--------|------|
|  1 | Setting_Update |
|  2 | User_Update |
|  3 | Channel_Create |
|  4 | Channel_Update |
|  5 | Channel_Delete |
|  6 | Role_Create |
|  7 | Role_Update |
|  8 | Role_Delete |
|  9 | Invite_Create |
| 10 | Invite_Update |
| 11 | Invite_Delete |
| 12 | Invite_Use |

# Remote Ok
Remote Ok Bot

## ❔ How it Works
```mermaid
flowchart LR
    A(RemoteOk API)
    B(Filter)
    C{{RemoteOk - Supabase}}
    D[(Supabase)]
    E(Send to Telegram)
    F(Save to Supabase)
    G[[Cleanup]]
    A --> B
    B --> C
    D --> C
    C --> E
    C --> F
```


1. `RemoteOk API` ➡️ Pull jobs posting from remoteOk API then applied filter criteria:
    - Job posting not older than 7 days

2. `Supabase` ➡️ Get All Job from Supabase Db

3. `RemoteOk - Supabase` ➡️ Keep RemoteOk Job Posting that is not in Supabase

4. `Send to Telegram` ➡️ Send Job(s) to Telegram Channel as a Bot

5. `Save to Supabase` ➡️ Save Job(s) to Supabase Db

**note:** _Job in the db will be cleanup(old job post) daily to reduce db size_
```
Worker run everyday at 10.00 AM
Cleaner run everday at 10.30 AM
```

## 🛠️ Setup
### Environment Variables
| Name                        | Desc                   |
|-----------------------------|------------------------|
| TELEGRAM_BOT                | Telegram bot API Token |
| TELEGRAM_CHANNEL_REMOTE_OK  | Telegram Channel id    |
| SUPABASE_URL                | Supabase URL           |
| SUPABASE_KEY                | Supabase Key           |

### Supabase Table Definition
```sql
create table
  public.RemoteOk (
    id text not null,
    epoch bigint not null,
    slug text not null default ''::text,
    company text not null default ''::text,
    position text not null default ''::text,
    description text not null default ''::text,
    location text not null default ''::text,
    url text not null default ''::text,
    constraint RemoteOk_pkey primary key (id)
  ) tablespace pg_default;
```
# Database Decision

## Recommendation

Do not use MongoDB as the long-term production database for this app.

Use:

- **Primary source of truth:** serverless Postgres, preferably Neon
- **Optional hot-path layer:** Upstash Redis

## Why

This is a Reddit-style application. The important data model is relational:

- users own posts and comments
- posts belong to communities
- comments belong to posts
- each user can vote once per post/comment
- saved items belong to users
- moderation/reporting/search will need queryable relationships later

MongoDB can work, but it adds less value here than a relational database. Redis alone is fast, but it makes correctness and future querying harder because uniqueness, foreign-key-like constraints, and joins become application-owned complexity.

## Selected architecture

```text
Angular frontend
  -> Go API
    -> Neon Postgres     canonical users/posts/comments/votes/saved data
    -> Upstash Redis     optional cache, rate limits, hot feed keys, temporary queues
```

Postgres owns durable truth. Redis can be deleted/rebuilt without losing user content.

## Why Neon Postgres first

Neon provides serverless Postgres with no server management, autoscaling, usage-based billing, and scale-to-zero behavior. That matches the goal of avoiding always-on database infrastructure while keeping SQL correctness.

## Where Upstash fits

Use Upstash Redis for:

- rate limiting
- session/token denylist if needed
- feed cache
- hot vote/comment counters cache
- short-lived queues or notification fanout

Do not make Redis the only source of truth for posts/comments/users unless this project intentionally trades correctness/queryability for speed.

## Candidate schema

```sql
create table users (
  id uuid primary key,
  username text not null unique,
  email text not null unique,
  password_hash text not null,
  created_at timestamptz not null default now()
);

create table communities (
  id uuid primary key,
  name text not null unique,
  created_by uuid references users(id),
  created_at timestamptz not null default now()
);

create table posts (
  id uuid primary key,
  community_id uuid not null references communities(id),
  user_id uuid not null references users(id),
  title text not null,
  body text,
  up_votes integer not null default 0,
  down_votes integer not null default 0,
  comments_count integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table comments (
  id uuid primary key,
  post_id uuid not null references posts(id),
  user_id uuid not null references users(id),
  body text not null,
  up_votes integer not null default 0,
  down_votes integer not null default 0,
  edited boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table votes (
  id uuid primary key,
  target_type text not null check (target_type in ('post', 'comment')),
  target_id uuid not null,
  user_id uuid not null references users(id),
  value smallint not null check (value in (-1, 1)),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique(target_type, target_id, user_id)
);

create table saved_items (
  id uuid primary key,
  user_id uuid not null references users(id),
  target_type text not null check (target_type in ('post', 'comment')),
  target_id uuid not null,
  created_at timestamptz not null default now(),
  unique(user_id, target_type, target_id)
);

create index posts_feed_idx on posts (created_at desc, id desc);
create index posts_community_feed_idx on posts (community_id, created_at desc, id desc);
create index comments_post_idx on comments (post_id, created_at asc, id asc);
create index votes_target_idx on votes (target_type, target_id);
```

## Migration plan

1. Add repository interfaces around persistence so handlers stop calling Mongo collections directly.
2. Add Postgres driver and migrations.
3. Implement users, posts, comments, votes, communities, and saved repositories in Postgres.
4. Replace `MONGOURI` with `DATABASE_URL`.
5. Keep cursor pagination, but change cursors from Mongo ObjectIDs to `(created_at, id)` based cursors.
6. Add Upstash Redis only after Postgres is working.
7. Use Redis for rate limits/cache only, not canonical content.

## Deployment env vars after migration

```bash
DATABASE_URL=postgres://...
SECRET_KEY=<long-random-secret>
ALLOWED_ORIGINS=https://<frontend-domain>
PORT=<set-by-host>

# optional
UPSTASH_REDIS_REST_URL=https://...
UPSTASH_REDIS_REST_TOKEN=...
```

## Success criteria

- No production dependency on MongoDB.
- DB enforces one vote per user per target.
- Feed and comments are still cursor-paginated.
- App can run with Postgres only.
- Redis can be cleared without losing user-generated content.

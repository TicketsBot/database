# Entitlements
The following entity relationship diagram describes the entities and relationships required to store both user and
server entitlements, with associated metadata, which are sourced from Patreon, Discord, voting and keys.

```mermaid
erDiagram
    wl[whitelabel_bot] {
        uuid entitlement_id
        bigint bot_id
        string token
        string public_key
        bool unrestricted_guild_limit
    }

    wlguilds[whitelabel_guilds] {
        bigint bot_id
        bigint guild_id
    }

    wlguilds }o--|| wl: associated

    uent[user_entitlements] {
        uuid entitlement_id
        bigint user_id
        uuid sku_id
        string source
        datetime expires_at
    }

    wl }o--|| uent : has
    uent ||--|| skus : associated
    gent ||--|| skus : associated
    uent ||--|| source : from
    gent ||--|| source : from

    gent[guild_entitlements] {
        uuid entitlement_id
        bigint guild_id
        bigint user_id
        uuid sku_id
        string source
        datetime expires_at
    }

    tiers {
        string label
        int priority
    }

    skus {
        uuid id
        string name
        string tier_label
    }

    tiers }o--|| skus : has

    wlskus[whitelabel_skus] {
        uuid sku_id
        int bots_permitted
        int servers_per_bot_permitted
    }

    wlskus |o--|| skus : links

    source[premium_sources] {
        string name
    }

    patreon[patreon_subscriptions] {
        bigint user_id
        uuid sku_id
        datetime expires_at
    }

    patreon ||--o{ tiers : has
    patreon ||--|| uent : links

    discord_store_skus {
        bigint discord_id
        uuid sku_id
    }

    discord_store_skus ||--o{ skus : links
```
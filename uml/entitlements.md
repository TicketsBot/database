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
    }

    wlguilds[whitelabel_guilds] {
        bigint bot_id
        bigint guild_id
    }

    wlguilds }o--|| wl: associated

    wl }o--o| ent : has
    ent ||--|| skus : associated

    ent[entitlements] {
        uuid entitlement_id
        bigint guild_id
        bigint user_id
        uuid sku_id
        enum_source source
        datetime expires_at
    }

    skus {
        uuid id
        string label
        string sku_type
    }

    skus |o--|| subscription_skus : associates

    subscription_skus {
        uuid sku_id
        string tier_label
        int priority
        bool is_global
    }

    wlskus[whitelabel_skus] {
        uuid sku_id
        int bots_permitted
        int servers_per_bot_permitted
    }

    wlskus |o--|| skus : links

    patreon[patreon_subscriptions] {
        bigint user_id
        uuid sku_id
        datetime expires_at
    }

    patreon ||--|| ent : links

    discord_store_skus {
        bigint discord_id
        uuid sku_id
    }

    discord_store_skus ||--o{ skus : links
```
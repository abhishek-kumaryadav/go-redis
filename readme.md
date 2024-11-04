


https://medium.com/@iggeehu/learning-go-by-writing-a-simple-tcp-server-d8ed260f67ac
https://app.codecrafters.io/courses/redis/introduction
https://www.honeybadger.io/blog/golang-logging/
https://redis.io/kb/doc/1fqjridk8w/what-are-the-impacts-of-the-redis-expiration-algorithm
https://medium.com/@harshilsokhadia/redis-scan-explained-an-efficient-way-to-extract-keys-3fd120f43b0e
https://www.linkedin.com/pulse/replication-redis-database-yeshwanth-n/

https://medium.com/@daijue/the-basics-of-replication-in-redis-4b92a3b275bd

https://stackoverflow.com/questions/50848042/what-is-internal-implementation-for-redis-master-slave-communication


go

https://emretanriverdi.medium.com/graceful-shutdown-in-go-c106fe1a99d9
https://emretanriverdi.medium.com/concurrency-in-go-1af9349b1103
https://github.com/huandu/skiplist
https://youtu.be/ol-FaNLXlR0?si=twZXb6XOV122feu2

https://pkg.go.dev/github.com/bigkevmcd/go-configparser#section-readme

internal/
├── auth/           # Authentication logic
├── database/       # Database interactions
├── middleware/     # HTTP middleware
├── models/         # Internal data structures
├── config/         # Configuration handling
├── validation/     # Input validation
└── service/        # Business logic

pkg/
├── api/           # API clients/interfaces
├── logger/        # Logging utilities
├── utils/         # Shared utilities
├── types/         # Shared type definitions
├── errors/        # Error definitions
└── client/        # Client libraries



Let me explain Redis's internal storage organization in detail.

Main Dictionary:
1. Primary data structure is a hash table where actual key-value pairs live
2. Each entry contains:
    - Key (string)
    - Value pointer (points to the actual value object)
    - Value type (string, list, hash, set, etc.)
    - Reference count (for memory management)
    - LRU information (for eviction)
    - Encoding information (how value is stored internally)

Metadata Dictionaries:
1. Expires Dictionary
    - Maps keys to their absolute Unix timestamp for expiration
    - Only contains entries for keys that have an expiry set
    - Checked during key access and by periodic cleanup task
    - When a key expires, it's removed from both expires and main dict

2. Blocking Information Dictionary
    - Used for BLPOP/BRPOP operations
    - Tracks which clients are waiting for which keys
    - Cleaned up when client disconnects or operation completes

3. Watch Dictionary
    - Used for WATCH command in transactions
    - Tracks which clients are watching which keys
    - Maps keys to list of clients watching them

Key type information:
- Type info stored with the key in main dictionary
- Common types: STRING, LIST, HASH, SET, ZSET
- Each type has different internal encodings:
    - STRING: int, embstr, raw
    - LIST: ziplist, linkedlist
    - HASH: ziplist, hashtable
    - SET: intset, hashtable
    - ZSET: ziplist, skiplist

Memory Management:
1. Redis uses a shared objects pool for:
    - Common integers (-128 to 127)
    - Common strings
    - Empty keys
2. Reference counting tracks object usage
3. Memory is freed when ref count reaches zero

Command Table:
- Separate hash table mapping command names to implementation
- Not visible to users through KEYS
- Includes command metadata like arity, flags

Would you like me to explain more about any specific aspect of these storage mechanisms?
---
layout: post
title: Enums
date: 2019-10-15
---

Enum (short for Enumeration) is GraphQL special type that represents set of predefined values. You'll mostly use this type when you need more options than just Boolean, but still want to use only controlled set of values.

```graphql
enum articleStatus {
  IDEA
  DRAFT
  PUBLISHED
}
```

## Resources

- [How to use GraphQL enum type and its best practices](https://graphqlmastery.com/blog/how-to-use-graphql-enum-type-and-its-best-practices)
- [Schemas and Types | GraphQL](https://graphql.org/learn/schema/#enumeration-types)
# API文档

## 网关路由说明（Nginx 同端口反代）

- Management 服务：保持原路径（例如 `/api/login`、`/api/admin/pools`）
- Game 服务：统一增加前缀 `/game`（例如 `/game/api/pull`）

对应后端端口：

- Management: `127.0.0.1:8080`
- Game: `127.0.0.1:8081`

## 健康检查

### Request

- Method: `GET`
- URL: `/ping`

### Response

- Body
  ```json
  {
    "code": 200,
    "message": "pong"
  }
  ```

---

## 注册

### Request

- Method: `POST`
- URL: `/api/register`
- Body:

  ```json
  {
    "username": "string(必须,最大40)",
    "nickname": "string(非必须，最大40)",
    "profile": "string(非必须，最大255)",
    "password": "string(必须,最大40)",
    "email": "string(必须，最大100)"
  }
  ```

### Response

- Body
  ```json
  {
    "code": 200,
    "uid": 100000000,
    "message": "注册成功",
    "role": "admin或user"
  }
  ```

---

## 登录

### Request

- Method: `POST`
- URL: `/api/login`
- Body:

  ```json
  {
    "username": "string(必须)",
    "uid": "uint(必须)",
    "email": "string(必须)",
    "password": "string(必须)"
  }
  ```

- 说明：当前实现会同时校验 `username`、`uid`、`email`、`password` 均非空。

### Response

- Body
  ```json
  {
    "code": 200,
    "message": "登录成功",
    "role": "admin或user",
    "access_token": "string",
    "refresh_token": "string"
  }
  ```

---

## 刷新令牌

### Request

- Method: `POST`
- URL: `/api/refresh`
- Body:

  ```json
  {
    "uid": 10001,
    "refresh_token": "string(必须)"
  }
  ```

### Response

- Body
  ```json
  {
    "code": 200,
    "message": "刷新成功",
    "access_token": "string",
    "refresh_token": "string"
  }
  ```

---

## 更新个人信息

### Request

- Method: `PUT`
- URL: `/api/auth/user/me`
- Auth: `Bearer 用户 access_token`
- Body:

  ```json
  {
    "nickname": "string(非必须，最大40)",
    "profile": "string(非必须，最大255)"
  }
  ```

### Response

- Body
  ```json
  {
    "code": 200,
    "message": "更新成功"
  }
  ```

---

## 抽卡

### Request

- Method: `GET`
- URL: `/game/api/pull`
- Auth: `Bearer 用户 access_token`
- Query:
  - `pool_id`: `uint(必须)`
  - `times`: `string(必须, 仅支持 "1" 或 "10")`

### Response

- 当 `times=1`:
  ```json
  {
    "id": 1,
    "name": "string",
    "rank": "S",
    "is_limited": true,
    "is_up": true,
    "code": 200,
    "message": "抽卡成功"
  }
  ```
- 当 `times=10`:
  ```json
  {
    "characters": [
      {
        "id": 1,
        "name": "string",
        "rank": "A",
        "is_limited": false,
        "is_up": false
      }
    ],
    "code": 200,
    "message": "抽卡成功"
  }
  ```

---

## 获取用户角色列表

### Request

- Method: `GET`
- URL: `/game/api/characters`
- Auth: `Bearer 用户 access_token`

### Response

- Body

  ```json
  {
    "uid": 10001,
    "characters": [
      {
        "character_id": 10,
        "character": {
          "id": 10,
          "name": "string",
          "rank": "A",
          "is_limited": false,
          "is_up": false
        },
        "owned_count": 2,
        "level": 1,
        "first_acquired_at": "2026-04-11T12:00:00Z",
        "first_acquired_pool": 1,
        "first_acquired_pull_count": 8
      }
    ],
    "code": 200,
    "message": "获取角色列表成功"
  }
  ```

- 说明：
  - 响应为 `uid` + `characters` + `code` + `message`。
  - 不再直接返回数据库模型对象。
  - `characters` 中每项不包含 `UserID` / `User` 等内部关联字段。

---

## 获取卡池ID列表

### Request

- Method: `GET`
- URL: `/game/api/pools`
- Auth: `Bearer 用户 access_token`

### Response

- Body
  ```json
  {
    "pool_ids": [1, 2, 3],
    "code": 200,
    "message": "获取卡池ID列表成功"
  }
  ```

---

## 获取卡池信息

### Request

- Method: `GET`
- URL: `/game/api/pool`
- Auth: `Bearer 用户 access_token`
- Query:
  - `pool_id`: `uint(必须)`

### Response

- Body

  ```json
  {
    "pool": {
      "id": 1,
      "name": "string",
      "description": "string",
      "Config": {
        "id": 1,
        "pool_id": 1,
        "s_rank_base_rate": 0.008,
        "a_rank_base_rate": 0.08,
        "a_guarantee_interval": 10,
        "s_pity_start": 65,
        "s_pity_step": 0.05,
        "s_pity_end": 80,
        "limit_pity": 120,
        "limit_rate_when_s": 0.5,
        "max_limited_characters": 1,
        "created_at": "2026-04-11T00:00:00Z",
        "updated_at": "2026-04-11T00:00:00Z"
      },
      "GachaPoolCharacters": [
        {
          "ID": 1,
          "PoolID": 1,
          "CharacterID": 10,
          "Character": {
            "id": 10,
            "name": "string",
            "rank": "S",
            "is_limited": true,
            "is_up": true
          },
          "CreatedAt": "2026-04-11T00:00:00Z"
        }
      ],
      "Characters": [
        {
          "id": 10,
          "name": "string",
          "rank": "S",
          "is_limited": true,
          "is_up": true
        }
      ],
      "start_at": "2026-04-11T00:00:00Z",
      "end_at": "2026-05-01T00:00:00Z",
      "UserPools": [
        {
          "ID": 1,
          "UserID": 10001,
          "PoolID": 1,
          "PullCount": 42,
          "LastACount": 40,
          "LastSCount": 35,
          "LastSUp": true,
          "LastUpCount": 35
        }
      ],
      "is_active": true
    },
    "code": 200,
    "message": "获取卡池信息成功"
  }
  ```

- 说明：
  - 保持原响应格式（`pool` + `code` + `message`）。
  - 不再直接返回数据库模型对象。
  - `pool` 返回卡池完整信息（`Config`、`GachaPoolCharacters`、`Characters`、`UserPools` 等字段）。

---

## 获取角色首次获取信息

### Request

- Method: `GET`
- URL: `/game/api/character/first_info`
- Auth: `Bearer 用户 access_token`
- Query:
  - `character_id`: `uint(必须)`

### Response

- Body
  ```json
  {
    "code": 200,
    "message": "获取角色首次信息成功",
    "first_acquired_at": "2026-04-11T12:00:00Z",
    "first_acquired_pool": 1,
    "first_acquired_pull_count": 37
  }
  ```

---

## 创建角色（管理员）

### Request

- Method: `POST`
- URL: `/api/admin/characters`
- Auth: `Bearer 管理员 access_token`
- Body:

  ```json
  {
    "name": "string(必须,最大40)",
    "rank": "string(必须, 例如 S/A/B)",
    "is_limited": false,
    "is_up": false
  }
  ```

### Response

- Body
  ```json
  {
    "id": 1,
    "name": "string",
    "rank": "S",
    "is_limited": false,
    "is_up": false,
    "code": 200,
    "message": "创建角色成功"
  }
  ```

---

## 创建卡池（管理员）

### Request

- Method: `POST`
- URL: `/api/admin/pools`
- Auth: `Bearer 管理员 access_token`
- Body:

  ```json
  {
    "pool": {
      "name": "string(必须,最大100)",
      "description": "string(非必须，最大255)",
      "start_at": "2026-04-11T00:00:00Z",
      "end_at": "2026-05-01T00:00:00Z",
      "is_active": true
    },
    "config": {
      "s_rank_base_rate": 0.008,
      "a_rank_base_rate": 0.08,
      "a_guarantee_interval": 10,
      "s_pity_start": 65,
      "s_pity_step": 0.05,
      "s_pity_end": 80,
      "limit_pity": 120,
      "limit_rate_when_s": 0.5,
      "max_limited_characters": 1
    }
  }
  ```

### Response

- Body
  ```json
  {
    "code": 200,
    "pool_id": 1,
    "message": "创建卡池成功"
  }
  ```

---

## 向卡池插入角色（管理员）

### Request

- Method: `POST`
- URL: `/api/admin/pools/characters`
- Auth: `Bearer 管理员 access_token`
- Body:

  ```json
  {
    "pool_id": 1,
    "character_id": 10
  }
  ```

### Response

- Body
  ```json
  {
    "code": 200,
    "message": "插入角色成功"
  }
  ```

---

## 更新卡池配置（管理员）

### Request

- Method: `PUT`
- URL: `/api/admin/pools/config`
- Auth: `Bearer 管理员 access_token`
- Body:

  ```json
  {
    "pool_id": 1,
    "s_rank_base_rate": 0.008,
    "a_rank_base_rate": 0.08,
    "a_guarantee_interval": 10,
    "s_pity_start": 65,
    "s_pity_step": 0.05,
    "s_pity_end": 80,
    "limit_pity": 120,
    "limit_rate_when_s": 0.5,
    "max_limited_characters": 1
  }
  ```

### Response

- Body
  ```json
  {
    "code": 200,
    "pool_id": 1,
    "version": 5,
    "message": "更新卡池配置成功"
  }
  ```

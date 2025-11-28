# 后端API文档

## 基础信息

### 服务地址
- **网关地址**: `http://your-domain.com` (Nginx反向代理)
- **用户服务**: 内部端口 8081
- **视频服务**: 内部端口 8082
- **弹幕服务**: 内部端口 8083
- **交互服务**: 内部端口 8084

### 认证方式
使用JWT Bearer Token认证
```
Authorization: Bearer <token>
```

需要认证的接口会返回401错误如果token无效或缺失。

---

## 1. 用户服务 (User Service)

### 1.1 用户注册
- **URL**: `POST /api/v1/user/register`
- **认证**: 不需要
- **请求体**:
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```
- **响应**:
```json
{
  "user": {
    "id": 1,
    "username": "string",
    "email": "string",
    "created_at": "2025-11-24T00:00:00Z"
  }
}
```

### 1.2 用户登录
- **URL**: `POST /api/v1/user/login`
- **认证**: 不需要
- **请求体**:
```json
{
  "email": "string",
  "password": "string"
}
```
- **响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string"
  }
}
```

### 1.3 获取用户信息
- **URL**: `GET /api/v1/user/profile/:id`
- **认证**: 需要
- **响应**:
```json
{
  "user": {
    "id": 1,
    "username": "string",
    "email": "string",
    "avatar_url": "string",
    "created_at": "2025-11-24T00:00:00Z"
  }
}
```

---

## 2. 视频服务 (Video Service)

### 2.1 获取标签列表
- **URL**: `GET /api/v1/tags`
- **认证**: 不需要
- **响应**:
```json
{
  "tags": [
    {
      "id": 1,
      "name": "游戏",
      "created_at": "2025-11-24T00:00:00Z"
    },
    {
      "id": 2,
      "name": "音乐",
      "created_at": "2025-11-24T00:00:00Z"
    }
  ]
}
```

### 2.2 上传视频
- **URL**: `POST /api/v1/videos`
- **认证**: 需要 (从token获取user_id)
- **Content-Type**: `multipart/form-data`
- **请求参数**:
  - `video`: File (视频文件)
  - `title`: string (标题)
  - `description`: string (描述)
  - `tag_ids[]`: array of integers (标签ID数组，可多选)
    - 例如: `tag_ids[]=1&tag_ids[]=3`

- **响应**:
```json
{
  "video": {
    "id": 1,
    "title": "string",
    "description": "string",
    "url": "http://domain.com/uploads/video.mp4",
    "user_id": 1,
    "created_at": "2025-11-24T00:00:00Z"
  }
}
```

### 2.3 获取视频流 (Feed)
- **URL**: `GET /api/v1/videos/feed`
- **认证**: 不需要
- **查询参数**:
  - `limit`: integer (每页数量，默认20)
  - `offset`: integer (偏移量，默认0)
  - `tag_id`: integer (标签过滤，可选)

- **示例**: `GET /api/v1/videos/feed?limit=10&offset=0&tag_id=1`

- **响应**:
```json
{
  "videos": [
    {
      "id": 1,
      "title": "string",
      "description": "string",
      "url": "string",
      "user_id": 1,
      "view_count": 1000,
      "like_count": 100,
      "created_at": "2025-11-24T00:00:00Z"
    }
  ]
}
```

---

## 3. 弹幕服务 (Danmaku Service)

### 3.1 WebSocket连接
- **URL**: `ws://your-domain.com/api/v1/danmaku/ws/:video_id`
- **协议**: WebSocket
- **路径参数**: 
  - `video_id`: integer (视频ID)

### 3.2 发送弹幕
通过WebSocket连接发送JSON消息:
```json
{
  "content": "弹幕内容",
  "time": 12.5,
  "color": "#FFFFFF",
  "type": 1
}
```
字段说明:
- `content`: string (弹幕文字内容)
- `time`: float (视频时间点，秒)
- `color`: string (颜色，十六进制)
- `type`: integer (类型: 1=滚动, 2=顶部, 3=底部)

### 3.3 接收弹幕
WebSocket会推送该视频房间的所有弹幕消息，格式同上。

### 3.4 获取历史弹幕
- **URL**: `GET /api/v1/danmaku/:video_id`
- **认证**: 不需要
- **路径参数**: `video_id` (视频ID)
- **响应**:
```json
{
  "danmaku_list": [
    {
      "id": 1,
      "video_id": 1,
      "user_id": 1,
      "content": "弹幕内容",
      "time": 12.5,
      "color": "#FFFFFF",
      "type": 1,
      "created_at": "2025-11-24T00:00:00Z"
    }
  ]
}
```

---

## 4. 交互服务 (Interaction Service)

### 4.1 点赞视频
- **URL**: `POST /api/v1/interaction/like`
- **认证**: 需要
- **请求体**:
```json
{
  "video_id": 1
}
```
- **响应**:
```json
{
  "message": "liked"
}
```

### 4.2 取消点赞
- **URL**: `DELETE /api/v1/interaction/like`
- **认证**: 需要
- **请求体**:
```json
{
  "video_id": 1
}
```
- **响应**:
```json
{
  "message": "unliked"
}
```

### 4.3 创建收藏夹
- **URL**: `POST /api/v1/interaction/favorites/folders`
- **认证**: 需要
- **请求体**:
```json
{
  "name": "收藏夹名称"
}
```
- **响应**:
```json
{
  "message": "folder created"
}
```

### 4.4 获取收藏夹列表
- **URL**: `GET /api/v1/interaction/favorites/folders`
- **认证**: 需要
- **响应**:
```json
{
  "folders": [
    {
      "id": 1,
      "user_id": 1,
      "name": "默认收藏夹",
      "created_at": "2025-11-24T00:00:00Z"
    }
  ]
}
```

### 4.5 添加视频到收藏夹
- **URL**: `POST /api/v1/interaction/favorites/add`
- **认证**: 需要
- **请求体**:
```json
{
  "folder_id": 1,
  "video_id": 1
}
```
- **响应**:
```json
{
  "message": "added to favorites"
}
```

---

## 5. 错误响应格式

所有错误响应统一格式:
```json
{
  "error": "错误描述信息"
}
```

常见HTTP状态码:
- `200`: 成功
- `400`: 请求参数错误
- `401`: 未认证或token无效
- `404`: 资源不存在
- `500`: 服务器内部错误

---

## 6. 数据库模型

### User (用户)
```go
type User struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"` // 不返回给前端
    AvatarURL string    `json:"avatar_url"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Video (视频)
```go
type Video struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    URL         string    `json:"url"`
    UserID      uint      `json:"user_id"`
    ViewCount   int       `json:"view_count"`
    LikeCount   int       `json:"like_count"`
    TagIDs      []uint    `json:"tag_ids" gorm:"-"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Tag (标签)
```go
type Tag struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Danmaku (弹幕)
```go
type Danmaku struct {
    ID        uint      `json:"id"`
    VideoID   uint      `json:"video_id"`
    UserID    uint      `json:"user_id"`
    Content   string    `json:"content"`
    Time      float64   `json:"time"`    // 视频时间点(秒)
    Color     string    `json:"color"`   // 十六进制颜色
    Type      int       `json:"type"`    // 1=滚动 2=顶部 3=底部
    CreatedAt time.Time `json:"created_at"`
}
```

---

## 7. 技术栈

- **语言**: Go 1.22
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 6.2
- **消息队列**: RabbitMQ 3.12
- **WebSocket**: gorilla/websocket
- **JWT**: golang-jwt/jwt

---

## 8. 开发建议

### 8.1 前端技术选型建议
- **框架**: React / Vue 3 / Next.js
- **视频播放器**: Video.js / DPlayer (支持弹幕)
- **WebSocket**: 原生WebSocket API / socket.io-client
- **HTTP客户端**: Axios / Fetch API
- **状态管理**: Redux / Vuex / Zustand

### 8.2 关键功能实现提示

#### JWT Token管理
```javascript
// 登录后保存token
localStorage.setItem('token', response.data.token);

// 请求时添加header
axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
```

#### WebSocket弹幕连接
```javascript
const ws = new WebSocket(`ws://your-domain.com/api/v1/danmaku/ws/${videoId}`);

ws.onmessage = (event) => {
  const danmaku = JSON.parse(event.data);
  // 在视频上渲染弹幕
};

// 发送弹幕
ws.send(JSON.stringify({
  content: "666",
  time: videoPlayer.currentTime,
  color: "#FFFFFF",
  type: 1
}));
```

#### 视频上传
```javascript
const formData = new FormData();
formData.append('video', videoFile);
formData.append('title', '视频标题');
formData.append('description', '描述');
formData.append('tag_ids[]', 1);
formData.append('tag_ids[]', 3);

axios.post('/api/v1/videos', formData, {
  headers: {
    'Content-Type': 'multipart/form-data',
    'Authorization': `Bearer ${token}`
  }
});
```

#### 分页加载
```javascript
const fetchVideos = async (page = 0, limit = 20, tagId = null) => {
  const params = { limit, offset: page * limit };
  if (tagId) params.tag_id = tagId;
  
  const response = await axios.get('/api/v1/videos/feed', { params });
  return response.data.videos;
};
```

---

## 9. 环境变量 (仅供参考)

后端使用的环境变量:
```bash
MYSQL_DSN=root:root@tcp(mysql:3306)/bilibili?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=redis:6379
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
JWT_SECRET=your-secret-key
STORAGE_PATH=./uploads
DOMAIN_URL=http://your-domain.com/uploads
```

---

## 10. CORS配置

后端已启用CORS，允许的请求头:
- `Content-Type`
- `Authorization`
- `X-Requested-With`

允许的方法: `GET, POST, PUT, DELETE, OPTIONS`

---

## 11. 联系方式

如有API问题或需要调整，请联系后端开发团队。

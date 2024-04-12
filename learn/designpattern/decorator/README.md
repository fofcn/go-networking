# 场景

我们应用开发经常会遇到认证和授权问题，比如：ERP系统、OA系统、CRM系统等等，这些系统都需要用户登录后才能访问，如何实现用户登录和权限认证呢？

我们来看下对应的解决方案：

## Python的装饰器模式
熟悉Python的大小朋友可能都知道Flask-JWT这个flask的extension,如何使用的？

```python
@jwt_required()
@app.route('/user/list')
def user_list():
    return users

```

## Java中认证授权
1. 放到应用Gateway去做
2. 自定义注解，类似于Python的Decorator
3. 自定义拦截器

## .NET的方式

```csharp

[Authorize(Roles = "Admin")]
public class UserController : Controller 
{
    [Route("api/users")]
    public IActionResult Get()
    {
        return Ok(new { message = "Hello World" });
    }
}

```

# 问题
1. 可能会忘记使用注解或者把他们放错位置
2. 阅读业务代码需要理解注解的意思
3. 找不到注解的定义
4. 控制流程被隐藏


# 解决方案

## 原生Go Http Server
Go的解决方案也比较简单，

首先定义一个认证包装器
```go

func Authentication(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 校验
        if !IsAuthenticated(r) {
            http.Error(w, "Not authenticated", http.StatusForbidden)
            return
        }
        h.ServeHTTP(w, r)
    }
}

```

然后使用它：
```go

Authentication(http.HandlerFunc(w http.ResponseWriter,r *http.Response){
    // 处理请求
})

```

最后注册路由：
```go

serve := http.NewServeMux()

// 需要验证的路由
serve.Handle("/",Authentication(http.HandlerFunc(w http.ResponseWriter,r *http.Response){
}))

// 不需要验证的路由
serve.Handle("/login",http.HandlerFunc(w http.ResponseWriter,r *http.Response){
})

```

## 基于Gin的认证
也是基于中间件的认证，和上面一样，只是使用gin框架。
也是首先定义一个Authentication中间件，然后注册路由的时候使用AuthMiddleware中间件。

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, common.CommonResp{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, common.CommonResp{
				Message: "Unauthorized",
			})
			c.Abort()
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			// 添加claims到上下文
			c.Set("claims", claims)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, common.CommonResp{
				Message: "Unauthorized",
			})
			c.Abort()
		}

	}
}

```
然后定义一个Public和protect group来分别处理公共接口和需要登录的接口。
```go

v1 := r.Group("/api/v1")

publicGroup := v1.Group("/")

protectGroup := v1.Group("/")
protectGroup.Use(user.AuthMiddleware())

user.InitRouter(publicGroup, protectGroup)

```


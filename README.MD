# 深拷贝

用法
```go

    type S struct {
        A int
        B int64
        C string
        inter
        T time.Time
    }
	
    type inter struct {
        loc string
    }
	
    func TestS(t *testing.T) {
        t.Log(Copy(S{
            A: 1,
            B: 2,
            C: "3",
            inter: inter{
                loc: "beijing",
            }, 
            T: time.Now(),
        }))
    }
```
输出：
```go
    {1 2 3 {} 2022-08-07 16:00:10.905688 +0800 CST m=+0.000465228}
```

注意：除了 time.Time 外，其他结构体的非导出字段不会复制

或者可以实现 DeepCopy Interface ，则可以实现复制结构体的非导出字段
```go
    type DeepCopy interface {
        DeepCopy() interface{}
    }
```

示例：
```go

    type dp struct {
        a int
    }
    
    func (d dp) DeepCopy() interface{} {
        return dp{d.a}
    }
    func TestDpInter(t *testing.T) {
        d := dp{4}
        t.Log(Copy(d))
    }

```

输出：
```go
    {4}
```